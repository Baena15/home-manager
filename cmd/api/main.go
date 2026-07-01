// ─── API Entry Point ──────────────────────
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gentleman-programming/home-manager/internal/config"
	"github.com/gentleman-programming/home-manager/internal/handlers"
	appmiddleware "github.com/gentleman-programming/home-manager/internal/middleware"
	"github.com/gentleman-programming/home-manager/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}

	db, err := store.New(cfg.DatabaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "database error: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if os.Getenv("AUTO_MIGRATE") == "true" || cfg.IsDevelopment() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := db.RunMigrations(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "migration error: %v\n", err)
			os.Exit(1)
		}
	}

	userStore := store.NewUserStore(db)

	if cfg.IsDevelopment() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := store.SeedUsers(ctx, userStore, cfg.OwnerEmail, cfg.OwnerPassword, cfg.PartnerEmail, cfg.PartnerPassword); err != nil {
			fmt.Fprintf(os.Stderr, "seed error: %v\n", err)
			os.Exit(1)
		}
	}

	authHandler := handlers.NewAuthHandler(cfg, userStore)
	authMiddleware := appmiddleware.AuthMiddleware(cfg.JWTSecret)

	productStore := store.NewProductStore(db)
	priceStore := store.NewPriceStore(db)
	productHandler := handlers.NewProductHandler(productStore, priceStore)

	listStore := store.NewListStore(db)
	listItemStore := store.NewListItemStore(db)
	listHandler := handlers.NewListHandler(listStore, listItemStore, productStore)

	expenseStore := store.NewExpenseStore(db)
	expenseHandler := handlers.NewExpenseHandler(expenseStore)

	incomeStore := store.NewIncomeStore(db)
	incomeHandler := handlers.NewIncomeHandler(incomeStore)

	settlementStore := store.NewSettlementStore(db)
	settlementHandler := handlers.NewSettlementHandler(settlementStore, userStore)

	dashboardStore := store.NewDashboardStore(db)
	dashboardHandler := handlers.NewDashboardHandler(dashboardStore)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Post("/api/v1/auth/login", authHandler.Login)
	r.Post("/api/v1/auth/register", authHandler.Register)

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.SetHeader("Content-Type", "application/json"))
		r.Use(authMiddleware)

		r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
			claims, ok := appmiddleware.ClaimsFromContext(r.Context())
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":"unauthorized"}`))
				return
			}
			respondJSON(w, http.StatusOK, claims)
		})

		r.Route("/products", func(r chi.Router) {
			r.Post("/", productHandler.Create)
			r.Get("/", productHandler.List)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", productHandler.Get)
				r.Put("/", productHandler.Update)
				r.Delete("/", productHandler.Delete)
				r.Post("/prices", productHandler.CreatePrice)
				r.Get("/prices", productHandler.ListPrices)
			})
		})

		r.Route("/lists", func(r chi.Router) {
			r.Post("/", listHandler.Create)
			r.Get("/", listHandler.List)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", listHandler.Get)
				r.Put("/", listHandler.Update)
				r.Delete("/", listHandler.Delete)
				r.Patch("/status", listHandler.UpdateStatus)
				r.Post("/items", listHandler.AddItem)
				r.Route("/items/{item_id}", func(r chi.Router) {
					r.Patch("/", listHandler.UpdateItem)
					r.Delete("/", listHandler.RemoveItem)
				})
			})
		})

		r.Route("/expenses", func(r chi.Router) {
			r.Post("/", expenseHandler.Create)
			r.Get("/", expenseHandler.List)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", expenseHandler.Get)
				r.Put("/", expenseHandler.Update)
				r.Delete("/", expenseHandler.Delete)
				r.Patch("/settle", expenseHandler.Settle)
			})
		})

		r.Route("/incomes", func(r chi.Router) {
			r.Post("/", incomeHandler.Create)
			r.Get("/", incomeHandler.List)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", incomeHandler.Get)
				r.Put("/", incomeHandler.Update)
				r.Delete("/", incomeHandler.Delete)
			})
		})

		r.Route("/settlements", func(r chi.Router) {
			r.Post("/", settlementHandler.Create)
			r.Get("/", settlementHandler.List)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", settlementHandler.Get)
				r.Delete("/", settlementHandler.Delete)
			})
		})

		r.Route("/dashboard", func(r chi.Router) {
			r.Get("/summary", dashboardHandler.Summary)
			r.Get("/monthly", dashboardHandler.Monthly)
			r.Get("/balance", dashboardHandler.Balance)
		})
	})

	fs := http.FileServer(http.Dir("./web"))
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})

	fmt.Printf("🚀 Server running on http://localhost:%s\n", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, "%v", payload)
}
