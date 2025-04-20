# 🧠 ToadML - Clasificación de Expresiones Faciales

Proyecto del módulo de MLOps para ITESO. Esta aplicación permite a los usuarios iniciar sesión, subir imágenes para realizar inferencias sobre expresiones faciales, visualizar resultados y consultar un historial de inferencias realizadas.

## 🌐 Descripción

La aplicación cuenta con un frontend atractivo tipo SPA (Single Page Application) desarrollado en **React**. Permite:

- Autenticación simple por nombre de usuario
- Subida de imágenes para inferencia
- Visualización de clase predicha y nivel de precisión (accuracy)
- Historial de inferencias con opción de descarga de imágenes
- Sección informativa del equipo

## 🖼️ Interfaz

- `Login`: Autenticación básica
- `Inferencia`: Carga de imagen, resultado de clase y precisión
- `Historial`: Tabla con inferencias previas, opción para descargar cada imagen
- `About`: Información del equipo y recursos utilizados

---

## 🚀 Tecnologías Utilizadas

### Frontend

- React
- Vite
- JavaScript
- CSS
- react-router-dom (SPA con navegación por scroll)
- HTML5 / CSS3

### Backend (esperado)

- FastAPI
- Modelo de clasificación en contenedor Docker
- Base de datos para guardar imágenes y resultados
- Autenticación, balanceo de carga, y nodos master/slave

### MLOps (esperado)

- MLFlow para experimentación
- DVC para versionado de datasets
- Git + GitHub para CI/CD
- Accuracy > 80%
- Early Stopping y curvas de aprendizaje

---

## 🛠️ Instalación y ejecución local

```bash
# Clonar el repositorio
git clone https://github.com/tuusuario/ToadML-FrontEnd.git
cd ToadML-FrontEnd/toad-mlops

# Instalar dependencias
npm install

# Ejecutar en desarrollo
npm run dev
```

## 📁 Estructura del Proyecto
```bash
toad-mlops/
├── public/
├── src/
│   ├── components/
│   │   ├── Navbar.jsx
│   │   ├── Inference.jsx
│   │   ├── History.jsx
│   │   └── About.jsx
│   ├── pages/
│   │   ├── Login.jsx
│   │   └── MainPage.jsx
│   ├── styles/
│   │   └── *.css
│   └── App.jsx
├── package.json
└── README.md
```

## 👥 Equipo

Nombre del equipo: Toad

### Integrantes:

- Jose Cota

- Horacio

- Uri

- Roberto

- Marian
