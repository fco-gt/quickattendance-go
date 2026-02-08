# Escalabilidad en Go - Guía de Mejores Prácticas

## ✅ Lo que YA está bien en tu arquitectura

### 1. **Separación en Capas**
```
Handler → Service → Repository → Database
```
✅ Permite escalar horizontalmente
✅ Fácil de testear
✅ Bajo acoplamiento

### 2. **DTOs**
✅ API estable independiente de la DB
✅ Versionado de API fácil
✅ Validación centralizada

### 3. **Context**
✅ Permite timeouts y cancelación
✅ Propagación de valores (user ID, trace ID, etc.)

---

## 🚀 Optimizaciones para Escalar

### 1. **Problema: 2 Queries Secuenciales**

**Código Actual:**
```go
// Query 1: Buscar por nombre
existingAgency, err := s.repo.GetByName(ctx, req.Name)
if existingAgency != nil {
    return nil, domain.ErrAgencyExists
}

// Query 2: Buscar por dominio
existingAgency, err = s.repo.GetByDomain(ctx, req.Domain)
if existingAgency != nil {
    return nil, domain.ErrAgencyExists
}
```

**Problema:** 2 round-trips a la base de datos = más latencia

**Solución 1: Query Única**
```go
// En repository/gorm_agency.go
func (r *AgencyRepo) ExistsByNameOrDomain(ctx context.Context, name, domain string) (bool, error) {
    var count int64
    err := r.db.WithContext(ctx).
        Model(&domain.Agency{}).
        Where("name = ? OR domain = ?", name, domain).
        Count(&count).Error
    return count > 0, err
}

// En service/agency_service.go
exists, err := s.repo.ExistsByNameOrDomain(ctx, req.Name, req.Domain)
if err != nil {
    return nil, err
}
if exists {
    return nil, domain.ErrAgencyExists
}
```

**Solución 2: Índice Único en DB (MEJOR)**
```go
// En domain/agency.go
type Agency struct {
    Name   string `gorm:"uniqueIndex;not null"`
    Domain string `gorm:"uniqueIndex;not null"`
}

// La DB rechazará duplicados automáticamente
// Manejas el error en el servicio
if err := s.repo.Create(ctx, agency); err != nil {
    if errors.Is(err, gorm.ErrDuplicatedKey) {
        return nil, domain.ErrAgencyExists
    }
    return nil, err
}
```

---

### 2. **Problema: Bcrypt es CPU-Intensivo**

**Código Actual:**
```go
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
```

**Problema:** Bcrypt bloquea el goroutine (puede tomar 100-300ms)

**Solución: Worker Pool para Hashing**
```go
// internal/service/password_hasher.go
type PasswordHasher struct {
    jobs    chan hashJob
    workers int
}

type hashJob struct {
    password string
    result   chan hashResult
}

type hashResult struct {
    hash string
    err  error
}

func NewPasswordHasher(workers int) *PasswordHasher {
    h := &PasswordHasher{
        jobs:    make(chan hashJob, 100),
        workers: workers,
    }
    
    // Iniciar workers
    for i := 0; i < workers; i++ {
        go h.worker()
    }
    
    return h
}

func (h *PasswordHasher) worker() {
    for job := range h.jobs {
        hash, err := bcrypt.GenerateFromPassword(
            []byte(job.password), 
            bcrypt.DefaultCost,
        )
        job.result <- hashResult{
            hash: string(hash),
            err:  err,
        }
    }
}

func (h *PasswordHasher) Hash(ctx context.Context, password string) (string, error) {
    result := make(chan hashResult, 1)
    
    select {
    case h.jobs <- hashJob{password: password, result: result}:
        select {
        case res := <-result:
            return res.hash, res.err
        case <-ctx.Done():
            return "", ctx.Err()
        }
    case <-ctx.Done():
        return "", ctx.Err()
    }
}
```

**Uso:**
```go
type AgencyService struct {
    repo   domain.AgencyRepo
    hasher *PasswordHasher
}

hashedPassword, err := s.hasher.Hash(ctx, req.Password)
```

---

### 3. **Problema: Conversión Manual DTO ↔ Domain**

**Código Actual:**
```go
return &dto.AgencyResponse{
    ID:        agency.ID,
    Name:      agency.Name,
    Domain:    agency.Domain,
    Address:   agency.Address,
    Phone:     agency.Phone,
    IsActive:  agency.IsActive,
    CreatedAt: agency.CreatedAt,
    UpdatedAt: agency.UpdatedAt,
}
```

**Problema:** Propenso a errores, difícil de mantener

**Solución: Métodos de Conversión**
```go
// En internal/dto/agency_dto.go

// ToAgencyResponse convierte un modelo de dominio a DTO de respuesta
func ToAgencyResponse(agency *domain.Agency) *AgencyResponse {
    return &AgencyResponse{
        ID:        agency.ID,
        Name:      agency.Name,
        Domain:    agency.Domain,
        Address:   agency.Address,
        Phone:     agency.Phone,
        IsActive:  agency.IsActive,
        CreatedAt: agency.CreatedAt,
        UpdatedAt: agency.UpdatedAt,
    }
}

// ToAgency convierte un DTO de registro a modelo de dominio
func (r *RegisterAgencyRequest) ToAgency(passwordHash string) *domain.Agency {
    return &domain.Agency{
        ID:           uuid.New(),
        Name:         r.Name,
        Domain:       r.Domain,
        PasswordHash: passwordHash,
        Address:      r.Address,
        Phone:        r.Phone,
        IsActive:     true,
    }
}
```

**Uso:**
```go
// En el servicio
agency := req.ToAgency(string(hashedPassword))
if err := s.repo.Create(ctx, agency); err != nil {
    return nil, err
}
return dto.ToAgencyResponse(agency), nil
```

---

### 4. **Caching para Lecturas Frecuentes**

```go
// internal/service/agency_service.go
type AgencyService struct {
    repo  domain.AgencyRepo
    cache *redis.Client
}

func (s *AgencyService) GetByID(ctx context.Context, id string) (*dto.AgencyResponse, error) {
    // 1. Intentar desde cache
    cacheKey := fmt.Sprintf("agency:%s", id)
    cached, err := s.cache.Get(ctx, cacheKey).Result()
    if err == nil {
        var agency domain.Agency
        json.Unmarshal([]byte(cached), &agency)
        return dto.ToAgencyResponse(&agency), nil
    }
    
    // 2. Si no está en cache, buscar en DB
    agency, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // 3. Guardar en cache
    data, _ := json.Marshal(agency)
    s.cache.Set(ctx, cacheKey, data, 5*time.Minute)
    
    return dto.ToAgencyResponse(agency), nil
}
```

---

### 5. **Rate Limiting**

```go
// internal/transport/http/middleware/rate_limit.go
func RateLimitMiddleware(limiter *rate.Limiter) gin.HandlerFunc {
    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error": "Too many requests",
            })
            c.Abort()
            return
        }
        c.Next()
    }
}

// Uso en main.go
limiter := rate.NewLimiter(100, 200) // 100 req/s, burst de 200
router.Use(RateLimitMiddleware(limiter))
```

---

### 6. **Connection Pooling**

```go
// cmd/server/main.go
func setupDatabase() *gorm.DB {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }
    
    sqlDB, _ := db.DB()
    
    // Configuración para alta concurrencia
    sqlDB.SetMaxOpenConns(100)        // Máximo de conexiones abiertas
    sqlDB.SetMaxIdleConns(10)         // Conexiones idle
    sqlDB.SetConnMaxLifetime(time.Hour) // Tiempo de vida
    
    return db
}
```

---

### 7. **Graceful Shutdown**

```go
// cmd/server/main.go
func main() {
    router := setupRouter()
    
    srv := &http.Server{
        Addr:    ":8080",
        Handler: router,
    }
    
    // Iniciar servidor en goroutine
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen: %s\n", err)
        }
    }()
    
    // Esperar señal de interrupción
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    log.Println("Shutting down server...")
    
    // Timeout de 5 segundos para terminar requests en curso
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }
    
    log.Println("Server exiting")
}
```

---

### 8. **Observabilidad**

```go
// internal/middleware/logging.go
func LoggingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        duration := time.Since(start)
        log.Printf(
            "[%s] %s %s %d %v",
            c.Request.Method,
            c.Request.URL.Path,
            c.ClientIP(),
            c.Writer.Status(),
            duration,
        )
    }
}

// Métricas con Prometheus
var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration in seconds",
        },
        []string{"method", "endpoint"},
    )
)
```

---

## 📊 Checklist de Escalabilidad

### Base de Datos
- [ ] Índices en columnas de búsqueda frecuente
- [ ] Connection pooling configurado
- [ ] Queries optimizadas (evitar N+1)
- [ ] Paginación en listados
- [ ] Read replicas para lecturas

### Aplicación
- [ ] Context con timeouts
- [ ] Graceful shutdown
- [ ] Rate limiting
- [ ] Circuit breaker para servicios externos
- [ ] Worker pools para tareas CPU-intensivas

### Infraestructura
- [ ] Load balancer (Nginx, HAProxy)
- [ ] Múltiples instancias de la app
- [ ] Cache distribuido (Redis)
- [ ] Message queue para tareas async (RabbitMQ, Kafka)
- [ ] CDN para assets estáticos

### Monitoreo
- [ ] Logging estructurado
- [ ] Métricas (Prometheus)
- [ ] Tracing distribuido (Jaeger, OpenTelemetry)
- [ ] Alertas (PagerDuty, Slack)
- [ ] Health checks

---

## 🎯 Prioridades según Escala

### < 1,000 usuarios
- ✅ Tu arquitectura actual está bien
- Agrega: Logging, health checks, índices en DB

### 1,000 - 10,000 usuarios
- Agrega: Redis cache, connection pooling
- Optimiza: Queries, índices compuestos

### 10,000 - 100,000 usuarios
- Agrega: Load balancer, múltiples instancias
- Implementa: Rate limiting, circuit breakers

### > 100,000 usuarios
- Agrega: Message queues, microservicios
- Implementa: Sharding de DB, CQRS

---

## 💡 Conclusión

Tu arquitectura actual **SÍ es escalable** para la mayoría de casos de uso.

Las optimizaciones dependen de:
1. **Cuántos usuarios** esperas
2. **Qué operaciones** son más frecuentes
3. **Cuál es tu bottleneck** actual

**Regla de oro:** No optimices prematuramente. Mide primero, optimiza después.
