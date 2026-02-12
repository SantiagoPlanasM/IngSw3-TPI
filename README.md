# 🛍️ Sistema de Gestión de Pedidos

Sistema completo de gestión de pedidos con arquitectura de microservicios, construido con Go (Backend), React (Frontend) y MySQL (Base de Datos). Proyecto diseñado para implementar CI/CD con **GitHub Actions** .

## 📋 Tabla de Contenidos

- [Características](#características)
- [Stack Tecnológico](#stack-tecnológico)
- [Arquitectura](#arquitectura)
- [Requisitos Previos](#requisitos-previos)
- [Instalación](#instalación)
- [Uso](#uso)
- [Tests](#tests)
- [Docker](#docker)
- [CI/CD con GitHub Actions](#cicd-con-github-actions)
- [API Endpoints](#api-endpoints)

## ✨ Características

### Backend (API REST)
- ✅ Arquitectura en capas (Handlers, Services, Repositories)
- ✅ Lógica de negocio completa con validaciones
- ✅ Gestión de estados de pedidos (PENDING → CONFIRMED → SHIPPED / CANCELLED)
- ✅ Control automático de stock
- ✅ Unit tests con mocks
- ✅ Integration tests con base de datos real
- ✅ Health checks

### Frontend (React)
- ✅ Catálogo de productos
- ✅ Carrito de compras interactivo
- ✅ Creación de pedidos
- ✅ Historial de pedidos
- ✅ Gestión de estados (Confirmar, Enviar, Cancelar)
- ✅ UI responsive con Tailwind CSS

### DevOps
- ✅ Dockerfiles con multi-stage build
- ✅ Docker Compose para orquestación
- ✅ Configuración por variables de entorno
- ✅ Imágenes optimizadas
- ✅ Health checks

## 🚀 Stack Tecnológico

### Backend
- **Lenguaje:** Go 1.21
- **Framework:** Gin (Web Framework)
- **ORM:** GORM
- **Base de Datos:** MySQL 8.0
- **Testing:** Go testing package

### Frontend
- **Framework:** React 18
- **Build Tool:** Vite
- **Estilos:** Tailwind CSS
- **HTTP Client:** Axios
- **Servidor Producción:** Nginx

### DevOps
- **Contenedores:** Docker
- **Orquestación:** Docker Compose
- **CI/CD:** Railway¿?

## 🏗️ Arquitectura

```
┌─────────────────┐      ┌─────────────────┐      ┌─────────────────┐
│                 │      │                 │      │                 │
│   Frontend      │─────▶│   Backend API   │─────▶│     MySQL       │
│   (React)       │      │   (Go + Gin)    │      │                 │
│   Port: 80      │      │   Port: 8080    │      │   Port: 3306    │
│                 │      │                 │      │                 │
└─────────────────┘      └─────────────────┘      └─────────────────┘
```

### Estructura del Backend

```
backend/
├── cmd/api/main.go              # Punto de entrada
├── internal/
│   ├── domain/models.go         # Modelos de dominio
│   ├── handlers/                # Controladores HTTP
│   ├── services/                # Lógica de negocio
│   ├── repositories/            # Acceso a datos
│   └── config/database.go       # Configuración DB
└── tests/integration/           # Tests de integración
```

## 📦 Requisitos Previos

- Docker >= 20.10
- Docker Compose >= 2.0
- Go >= 1.21 (para desarrollo local)
- Node.js >= 18 (para desarrollo local)

## 🔧 Instalación

### 1. Clonar el repositorio

```bash
git clone <url-del-repositorio>
cd order-management-system
```

### 2. Configurar variables de entorno

```bash
cp .env.example .env
```

Editar `.env` según necesidad:

```env
MYSQL_ROOT_PASSWORD=rootpassword
MYSQL_DATABASE=order_management
MYSQL_USER=orderuser
MYSQL_PASSWORD=orderpass123
MYSQL_PORT=3306
BACKEND_PORT=8080
FRONTEND_PORT=80
```

### 3. Iniciar con Docker Compose

```bash
docker-compose up -d
```

Esto iniciará:
- MySQL en `localhost:3306`
- Backend API en `http://localhost:8080`
- Frontend en `http://localhost:80`

### 4. Verificar que todo esté funcionando

```bash
# Backend health check
curl http://localhost:8080/health

# Ver logs
docker-compose logs -f
```

## 💻 Uso

### Acceder a la aplicación

Abrir en el navegador: `http://localhost`

### Flujo de trabajo

1. **Ver Productos**: Tab "Productos" - Catálogo completo
2. **Agregar al Carrito**: Click en "Agregar al Carrito"
3. **Crear Pedido**: Tab "Carrito" - Seleccionar usuario y crear pedido
4. **Gestionar Pedidos**: Tab "Historial de Pedidos"
   - **Confirmar**: Reduce el stock (PENDING → CONFIRMED)
   - **Enviar**: Marca como enviado (CONFIRMED → SHIPPED)
   - **Cancelar**: Devuelve el stock si estaba confirmado

## 🧪 Tests

### Unit Tests (Backend)

```bash
cd backend
go test ./internal/services/... -v
```

### Integration Tests (Backend)

```bash
# Configurar variables de entorno para tests
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=orderuser
export DB_PASSWORD=orderpass123
export DB_NAME=order_management_test
export INTEGRATION_TEST=true

# Ejecutar tests
cd backend
go test ./tests/integration/... -v
```

### Ejecutar todos los tests

```bash
cd backend
go test ./... -v
```

## 🐳 Docker

### Construir imágenes manualmente

```bash
# Backend
docker build -t order-management-backend ./backend

# Frontend
docker build -t order-management-frontend ./frontend
```

### Comandos útiles

```bash
# Iniciar servicios
docker-compose up -d

# Ver logs
docker-compose logs -f

# Reiniciar un servicio
docker-compose restart backend

# Detener todo
docker-compose down

# Limpiar volúmenes (⚠️ elimina datos)
docker-compose down -v
```

## 🔄 CI/CD con GitHub Actions

### Pipeline Completo

El proyecto incluye un pipeline de CI/CD completo con GitHub Actions:

**Pipeline configurado en:** `.github/workflows/ci-cd.yml`

#### Stages del Pipeline:

1. **🧪 Unit Tests**
   - Ejecuta tests unitarios del backend
   - Genera reporte de cobertura
   - Sube resultados a Codecov
   - Comenta en PRs con resultados

2. **🔨 Build Images**
   - Construye imágenes Docker de Backend y Frontend
   - Pushea a GitHub Container Registry (ghcr.io)
   - Cachea layers para builds rápidos

3. **🚀 Deploy to QA**
   - Deploy automático a ambiente QA
   - Se ejecuta en pushes a `develop` o PRs
   - Health checks automáticos

4. **🔬 Integration Tests**
   - Ejecuta tests de integración en QA
   - Usa MySQL en GitHub Actions
   - Verifica funcionamiento completo

5. **🎯 Deploy to Production**
   - **Requiere aprobación manual**
   - Solo se ejecuta en rama `main`
   - Health checks y smoke tests
   - Crea GitHub Release automáticamente

## 📡 API Endpoints

### Users

```
GET    /api/users          # Listar todos los usuarios
GET    /api/users/:id      # Obtener usuario por ID
POST   /api/users          # Crear usuario
```

### Products

```
GET    /api/products       # Listar todos los productos
GET    /api/products/:id   # Obtener producto por ID
POST   /api/products       # Crear producto
```

### Orders

```
GET    /api/orders                 # Listar todos los pedidos
GET    /api/orders/:id             # Obtener pedido por ID
GET    /api/orders/user/:userId    # Obtener pedidos de un usuario
POST   /api/orders                 # Crear pedido
PATCH  /api/orders/:id/confirm     # Confirmar pedido
PATCH  /api/orders/:id/ship        # Enviar pedido
PATCH  /api/orders/:id/cancel      # Cancelar pedido
```

## 📝 Lógica de Negocio

### Estados de Pedido

```
PENDING ──┐
          ├─→ CONFIRMED ──→ SHIPPED
          └─→ CANCELLED
```

### Reglas

1. **PENDING → CONFIRMED**: Se valida y reduce el stock
2. **CONFIRMED → SHIPPED**: Solo se cambia el estado
3. **PENDING/CONFIRMED → CANCELLED**: Se devuelve el stock (si estaba confirmado)
4. **SHIPPED**: No se puede cancelar