# API de Gestión de Biblioteca

Este documento explica cómo construir y ejecutar la API de Gestión de Biblioteca usando Docker.

## 📋 Requisitos Previos

- Docker instalado (versión 20.10 o superior)
- Docker Compose instalado (opcional, para uso con docker-compose.yml)

## 🚀 Construcción de la Imagen

### Opción 1: Usando Docker directamente

```bash
# Construir la imagen
docker build -t books-api:latest .

# Verificar que la imagen se creó correctamente
docker images | grep books-api
```

### Opción 2: Usando Docker Compose

```bash
# Construir la imagen con docker-compose
docker-compose build
```

## 🏃 Ejecución del Contenedor

### Opción 1: Docker Run

#### Ejecución básica

```bash
docker run -d \
  --name books-api \
  -p 8080:8080 \
  books-api:latest
```

#### Ejecución con volúmenes persistentes

```bash
docker run -d \
  --name books-api \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/logs:/app/logs \
  books-api:latest
```

#### Ejecución con variables de entorno personalizadas

```bash
docker run -d \
  --name books-api \
  -p 8080:8080 \
  -e PORT=8080 \
  -e DB_PATH=/app/data/books.db \
  -e LOG_PATH=/app/logs/api.log \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/logs:/app/logs \
  books-api:latest
```

### Opción 2: Docker Compose

```bash
# Iniciar el servicio
docker-compose up -d

# Ver los logs en tiempo real
docker-compose logs -f

# Detener el servicio
docker-compose down

# Detener y eliminar volúmenes
docker-compose down -v
```

## 📊 Gestión del Contenedor

### Ver logs

```bash
# Docker
docker logs books-api

# Docker logs en tiempo real
docker logs -f books-api

# Docker Compose
docker-compose logs -f books-api
```

### Detener el contenedor

```bash
# Docker
docker stop books-api

# Docker Compose
docker-compose stop
```

### Reiniciar el contenedor

```bash
# Docker
docker restart books-api

# Docker Compose
docker-compose restart
```

### Eliminar el contenedor

```bash
# Docker
docker rm -f books-api

# Docker Compose
docker-compose down
```

### Acceder al contenedor

```bash
# Docker
docker exec -it books-api sh

# Docker Compose
docker-compose exec books-api sh
```

## 🔍 Verificación del Estado

### Verificar que el contenedor está corriendo

```bash
docker ps | grep books-api
```

### Verificar el healthcheck

```bash
docker inspect books-api --format='{{.State.Health.Status}}'
```

### Probar la API

```bash
# Desde el host
curl http://localhost:8080/configuration

# Desde dentro del contenedor
docker exec books-api wget -qO- http://localhost:8080/configuration
```

## 📁 Estructura de Volúmenes

La aplicación usa dos volúmenes persistentes:

- **`/app/data`**: Almacena la base de datos SQLite
  - Archivo: `books.db`
  - Mapeo sugerido: `./data:/app/data`

- **`/app/logs`**: Almacena los logs de la aplicación
  - Archivo: `api.log`
  - Mapeo sugerido: `./logs:/app/logs`

## 🌐 Endpoints Disponibles

Una vez que el contenedor esté corriendo, la API estará disponible en:

- Base URL: `http://localhost:8080`

Endpoints principales:

- `GET /authors` - Lista de autores
- `GET /books` - Lista de libros
- `GET /users` - Lista de usuarios
- `GET /loans` - Lista de préstamos
- `GET /reservations` - Lista de reservaciones
- `GET /fines` - Lista de multas
- Y muchos más...

## 🔧 Variables de Entorno

| Variable   | Descripción                       | Valor por defecto    |
| ---------- | --------------------------------- | -------------------- |
| `PORT`     | Puerto en el que escucha la API   | `8080`               |
| `DB_PATH`  | Ruta del archivo de base de datos | `/app/data/books.db` |
| `LOG_PATH` | Ruta del archivo de logs          | `/app/logs/api.log`  |

## 📦 Multi-Stage Build

El Dockerfile utiliza un build multi-stage para optimizar el tamaño de la imagen:

1. **Stage 1 (Builder)**: Compila la aplicación Go con todas las dependencias necesarias
2. **Stage 2 (Runtime)**: Crea una imagen ligera con solo el binario compilado

Ventajas:

- Imagen final más pequeña (~20-30 MB vs ~400+ MB)
- Mayor seguridad (menos superficie de ataque)
- Inicio más rápido del contenedor

## 🔒 Seguridad

El contenedor implementa las siguientes prácticas de seguridad:

- ✅ Ejecuta como usuario no-root (`appuser`)
- ✅ Imagen base Alpine (ligera y segura)
- ✅ Solo incluye dependencias de runtime necesarias
- ✅ Healthcheck configurado para monitoreo
- ✅ Puertos específicos expuestos

## 🐛 Troubleshooting

### El contenedor no inicia

```bash
# Ver logs detallados
docker logs books-api

# Verificar el estado del contenedor
docker inspect books-api
```

### No se pueden crear archivos en los volúmenes

```bash
# Verificar permisos de los directorios
ls -la data/ logs/

# Dar permisos si es necesario
chmod 755 data/ logs/
```

### La base de datos no persiste

```bash
# Verificar que los volúmenes están montados
docker inspect books-api --format='{{.Mounts}}'
```

### Error de compilación

```bash
# Limpiar la caché de Docker
docker builder prune

# Reconstruir sin caché
docker build --no-cache -t books-api:latest .
```

## 📝 Notas Adicionales

- La base de datos se crea automáticamente en el primer inicio
- Los logs se escriben tanto en consola como en archivo
- El healthcheck verifica el endpoint `/configuration` cada 30 segundos
- Se recomienda usar volúmenes para persistir datos en producción
