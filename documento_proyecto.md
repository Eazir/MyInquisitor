# MyInquisitor: Plataforma Modular de Control Financiero Personal

---

## Contenidos Asociados al Proyecto (Palabras Clave)

Control financiero, finanzas personales, gestión de deudas, contabilidad mensual, gastos recurrentes, flujo de caja, proyecciones financieras, arquitectura hexagonal, Go, Gin, React, TypeScript, PostgreSQL, SQLC, Docker, JWT, cifrado AES-256-GCM, componentes plantilla, TailwindCSS, código abierto, desacoplamiento, temas claro/oscuro, panel de administración.

---

## Resumen del Proyecto

MyInquisitor es una plataforma web de código abierto diseñada para el control financiero personal, construida sobre una arquitectura hexagonal en Go con el framework Gin, y un frontend basado en componentes plantilla con React, TypeScript y TailwindCSS. La plataforma permite a los usuarios gestionar deudas con seguimiento mes a mes, registrar gastos recurrentes, llevar una contabilidad mensual completa de ingresos y egresos, y visualizar balances y flujos de caja en periodos semanales, mensuales y anuales con proyecciones basadas en estimaciones del usuario. Incluye un sistema de autenticación robusto mediante tokens JWT con renovación automática, un panel de administración donde un super-usuario gestiona cuentas y permisos, y cifrado AES-256-GCM de todos los datos sensibles de los usuarios. El proyecto está completamente dockerizado en tres contenedores independientes (backend, frontend, base de datos PostgreSQL 16) y sigue un diseño de bajo acoplamiento donde cada funcionalidad puede modificarse o eliminarse sin afectar al resto. Los componentes del frontend operan como plantillas reutilizables que consumen variables CSS, permitiendo cambiar temas claro y oscuro, tipografías, tamaños y espaciados de forma global e inmediata.

---

## Descripción del Problema

### Contexto General

En la era digital actual, la gestión de las finanzas personales se ha convertido en una necesidad fundamental para individuos y hogares. La proliferación de servicios financieros digitales, suscripciones recurrentes, préstamos personales y modalidades de pago fraccionado ha incrementado significativamente la complejidad de mantener un control financiero efectivo. Según estudios recientes, más del 60% de las personas no lleva un registro sistemático de sus gastos mensuales, y una proporción aún mayor desconoce su flujo de caja real. Esta falta de control se traduce en endeudamiento no planificado, imposibilidad de ahorro y estrés financiero crónico.

### Problemas Identificados

#### 1. Fragmentación de herramientas financieras

El mercado actual ofrece una amplia variedad de aplicaciones financieras, pero ninguna aborda de manera integral las necesidades específicas del control financiero personal. Existen herramientas para llevar contabilidad empresarial (como QuickBooks o ContaAzul), aplicaciones para control de gastos personales (como Fintonic o YNAB), y soluciones bancarias nativas que ofrecen resúmenes limitados. Sin embargo, ninguna de estas soluciones integra de forma natural los siguientes aspectos:

- Seguimiento de deudas con estado de pago mes a mes.
- Gestión de gastos recurrentes (suscripciones, alquiler, servicios).
- Contabilidad mensual con ingresos y egresos detallados.
- Proyecciones y estimaciones basadas en datos históricos y entrada manual del usuario.
- Panel de control unificado con balances generales, disponibles y reservados.

El usuario se ve forzado a utilizar múltiples aplicaciones simultáneamente, lo que genera fricción, duplicación de esfuerzos y una visión fragmentada de su realidad financiera.

#### 2. Dificultad en el seguimiento de deudas a largo plazo

Las deudas personales requieren un seguimiento continuo que la mayoría de las herramientas no proporciona. El usuario necesita saber no solo el monto total pendiente, sino también el estado de pago de cada mes, el interés acumulado, y proyectar cuándo terminará de pagar. Las soluciones existentes tratan las deudas como un gasto único, perdiendo la granularidad necesaria para un control efectivo.

Una deuda típica tiene las siguientes dimensiones que deben ser rastreadas:

- Monto total original y monto restante.
- Cuota mensual y fecha de vencimiento.
- Estado de pago de cada mes (pagado, pendiente, atrasado).
- Intereses acumulados y proyección de liquidación.
- Capacidad de registrar pagos adicionales o adelantados.

Ninguna herramienta de propósito general cubre todas estas dimensiones de forma integrada con el resto de las finanzas del usuario.

#### 3. Falta de integración entre gastos recurrentes y contabilidad general

Los gastos recurrentes representan una porción significativa del presupuesto mensual. Suscripciones a servicios de streaming, alquiler, planes de telefonía, seguros, gimnasio, y otros gastos fijos pueden consumir entre el 40% y el 60% de los ingresos mensuales de una persona promedio. Sin embargo, rara vez se integran con la contabilidad general del usuario, lo que impide tener una visión holística de las finanzas.

El usuario termina consultando múltiples fuentes para entender su situación real:

- Estado de cuenta bancario para ver los cargos del mes.
- Aplicaciones de suscripciones para saber qué servicios tiene activos.
- Notas personales o spreadsheets para llevar la contabilidad.
- Recordatorios mentales o calendarios para fechas de pago.

Esta fragmentación cognitiva incrementa la probabilidad de olvidar pagos, incurrir en cargos por mora, y mantener suscripciones no utilizadas.

#### 4. Costo y privacidad de las soluciones comerciales

Las plataformas comerciales de control financiero presentan dos barreras importantes:

**Costo**: Las suscripciones mensuales o anuales pueden ser prohibitivas. YNAB, por ejemplo, cobra aproximadamente $100 USD al año. Otras plataformas similares tienen modelos freemium que limitan funcionalidades esenciales precisamente cuando el usuario más las necesita.

**Privacidad**: Estas plataformas almacenan los datos financieros de los usuarios en servidores externos, a menudo sin garantías suficientes de cifrado y transparencia sobre el procesamiento de datos. Dado que la información financiera es uno de los datos más sensibles que una persona puede compartir —revela ingresos, hábitos de consumo, capacidad de pago, deudas y patrimonio— la falta de control sobre su almacenamiento representa un riesgo significativo.

El usuario debe confiar ciegamente en que la empresa:
- No compartirá sus datos con terceros.
- Implementa medidas de seguridad adecuadas.
- No utilizará sus datos para perfiles comerciales.
- Mantendrá el servicio operativo a largo plazo.

#### 5. Acoplamiento fuerte en soluciones open source existentes

Las alternativas de código abierto disponibles adolecen de un acoplamiento fuerte entre sus módulos. Esto significa que modificar o extender una funcionalidad específica —como añadir un nuevo tipo de visualización, cambiar la lógica de proyecciones, o agregar una nueva fuente de datos— implica riesgos de afectar otras partes del sistema.

Las consecuencias de este acoplamiento incluyen:

- **Dificultad de mantenimiento**: Un cambio en un módulo requiere verificar todos los demás módulos.
- **Barrera de entrada para contribuciones**: Los nuevos desarrolladores deben entender todo el sistema antes de poder modificar una parte.
- **Riesgo de regresiones**: Cambios pequeños pueden tener efectos colaterales imprevistos.
- **Dificultad de evolución**: Adaptarse a nuevos requisitos es lento y costoso.

#### 6. Limitaciones en la personalización visual

La mayoría de las aplicaciones financieras ofrecen temas visuales fijos o limitados a una o dos variantes predefinidas. Cambiar la apariencia global —colores, tipografías, espaciados, radios de borde— requiere modificaciones caso por caso en cada pantalla o componente.

Esto tiene implicaciones directas:

- **Esfuerzo de mantenimiento**: Cambiar un color primario implica editar decenas de archivos.
- **Accesibilidad**: Adaptar la interfaz para usuarios con discapacidades visuales requiere cambios generalizados.
- **Experiencia de usuario**: No se puede ofrecer personalización visual sin un costo de desarrollo significativo.

### Solución Propuesta

MyInquisitor aborda estos seis problemas de manera integral mediante un diseño arquitectónico cuidadoso y decisiones tecnológicas específicas que se describen a continuación.

#### Arquitectura hexagonal para desacoplamiento total

La adopción de una arquitectura hexagonal (también conocida como puertos y adaptadores) garantiza que la lógica de negocio esté completamente aislada de los detalles de infraestructura. El dominio y los casos de uso no dependen de frameworks, bases de datos o protocolos de comunicación específicos. Esto se traduce en:

- Cada funcionalidad (deudas, gastos, contabilidad, administración) es un módulo independiente con sus propias interfaces y casos de uso.
- Cambiar la base de datos, el framework HTTP o añadir una interfaz móvil no requiere modificar la lógica de negocio.
- El testing se simplifica al poder mockear las interfaces de repositorio sin necesidad de infraestructura real.
- El sistema puede crecer orgánicamente añadiendo nuevos módulos sin riesgos de regresiones.

#### Componentes plantilla para consistencia visual global

El frontend se construye con componentes React que funcionan como plantillas reutilizables. Estos componentes no contienen valores concretos de color, tamaño o tipografía, sino que referencian variables CSS definidas a nivel global. Este enfoque permite:

- Cambiar toda la apariencia del sistema modificando únicamente las variables CSS en un archivo central.
- Implementar temas claro y oscuro completos sin tocar una sola página o componente individual.
- Ajustar la responsividad, accesibilidad y espaciado desde un punto central de configuración.
- Reutilizar los mismos componentes en una futura aplicación móvil construida con React Native, manteniendo la consistencia visual entre plataformas.

#### Seguridad desde el diseño

La plataforma incorpora múltiples capas de seguridad que protegen los datos del usuario en cada etapa:

- **Autenticación**: JWT con tokens de acceso de corta duración (15 minutos) y tokens de refresco de larga duración (7 días), con renovación automática y transparente para el usuario.
- **Cifrado en reposo**: AES-256-GCM para todos los datos personales identificables (email, nombre completo, teléfono). La clave de cifrado es única por entorno y se suministra como variable de entorno.
- **Hashing de contraseñas**: bcrypt con costo computacional 12, que proporciona un equilibrio óptimo entre seguridad y rendimiento.
- **Protección contra inyección SQL**: SQLC genera código Go con queries parametrizadas, eliminando por completo este vector de ataque.
- **Validación de entrada**: Todos los endpoints validan el cuerpo y los parámetros de la petición mediante go-playground/validator.

#### Arquitectura dockerizada para despliegue consistente

La plataforma se distribuye en tres contenedores Docker independientes, cada uno con su propio ciclo de vida, permitiendo despliegues consistentes en cualquier entorno:

- **Backend**: Aplicación Go compilada de forma multi-stage. La imagen de compilación (golang:1.23-alpine) contiene todas las herramientas de desarrollo, mientras que la imagen final (distroless) incluye únicamente el binario compilado, reduciendo el tamaño de ~1.2 GB a ~15 MB y minimizando la superficie de ataque.
- **Frontend**: Build estático generado con Vite y servido con Nginx. La imagen final solo contiene los archivos HTML, CSS y JS optimizados.
- **Base de datos**: PostgreSQL 16 con volumen persistente y health checks configurables.

---

## Objetivos del Proyecto

### Objetivo General

Desarrollar una plataforma web modular y desacoplada para el control financiero personal, implementada con arquitectura hexagonal en el backend y componentes plantilla en el frontend, que permita a los usuarios gestionar deudas, gastos recurrentes y contabilidad mensual de forma integrada, segura y personalizable.

### Objetivos Específicos

1. **Diseñar e implementar una base de datos PostgreSQL normalizada** que soporte las entidades necesarias para el control financiero personal: usuarios, deudas con seguimiento mensual, gastos recurrentes, transacciones contables, categorías y resúmenes mensuales precalculados, con migraciones versionadas y queries optimizadas mediante SQLC.

2. **Implementar un backend con arquitectura hexagonal en Go y Gin** que separe claramente las capas de dominio, aplicación e infraestructura, permitiendo que cada funcionalidad (autenticación, deudas, gastos, contabilidad, administración) sea independiente, testeable de forma aislada y modificable sin afectar al resto.

3. **Desarrollar un sistema de autenticación robusto basado en JWT** con tokens de acceso (15 minutos) y refresco (7 días), incluyendo renovación automática desde el frontend, protección de rutas según el rol del usuario, y almacenamiento seguro de tokens.

4. **Implementar cifrado AES-256-GCM para datos sensibles** garantizando que la información personal identificable (email, nombre, teléfono) se almacene cifrada en la base de datos y solo se descifre en memoria durante el procesamiento autorizado, con clave de cifrado configurable por entorno.

5. **Desarrollar un frontend con React, TypeScript y TailwindCSS basado en componentes plantilla** que permita la personalización global de temas (claro/oscuro), tipografías y espaciados mediante variables CSS, garantizando consistencia visual en toda la aplicación y facilitando la futura creación de una versión móvil.

6. **Crear un panel de administración** donde un super-usuario pueda gestionar cuentas de usuario (crear, modificar, activar/desactivar, cambiar permisos) de forma segura, con verificación de contraseña del usuario objetivo cuando sea necesario.

7. **Implementar dashboards y visualizaciones** que muestren balances generales (dinero total manejado, disponible, reservado), flujos de caja por periodos (semanal, mensual, anual), proyecciones financieras y estimaciones basadas en datos históricos e ingresos informados por el usuario.

8. **Dockerizar la plataforma completa** con contenedores independientes para backend, frontend y base de datos, utilizando builds multi-stage para optimizar el tamaño de las imágenes de producción y facilitar el despliegue reproducible en cualquier entorno.

9. **Garantizar un bajo acoplamiento entre funcionalidades** mediante la arquitectura hexagonal y las interfaces de repositorio en la capa de dominio, de modo que cada módulo pueda ser modificado, extendido o eliminado de forma independiente.

---

## Instrumentos Aplicados y/o Actividades Realizadas

### Análisis y Diseño

1. **Diseño de base de datos relacional**: Se diseñó un esquema normalizado de 12 tablas con relaciones bien definidas, claves foráneas, índices en columnas de búsqueda frecuente, y constraints de integridad referencial. Las tablas cubren usuarios, deudas, seguimiento mensual de deudas, gastos recurrentes, seguimiento mensual de gastos, transacciones contables, categorías y resúmenes mensuales precalculados.

2. **Diseño de arquitectura hexagonal**: Se definieron las capas del backend con reglas estrictas de dependencia unidireccional: la capa de dominio no conoce nada fuera de la stdlib de Go, la capa de aplicación solo conoce el dominio, y la infraestructura implementa las interfaces definidas en el dominio. Se establecieron los puertos (interfaces de repositorio) y adaptadores (implementaciones concretas con SQLC, JWT, Gin).

3. **Diseño de componentes plantilla**: Se catalogaron y diseñaron 15 componentes frontend reutilizables con interfaz tipada completa, soporte de variantes visuales, y dependencia exclusiva de variables CSS para todos los valores de estilo. Cada componente incluye estados de carga, vacío, error y éxito donde corresponde.

### Desarrollo

4. **Generación automática de código con SQLC**: Se escribieron queries SQL parametrizadas para cada entidad del sistema, cubriendo operaciones CRUD, búsquedas por usuario, consultas por periodo, y agregaciones contables. SQLC generó código Go tipado con interfaces, eliminando riesgos de inyección SQL y reduciendo el boilerplate manual en aproximadamente un 60%.

5. **Implementación de API REST con Gin**: Se desarrollaron 25+ endpoints HTTP organizados en grupos lógicos bajo el prefijo `/api/v1/`, con validación de entrada mediante go-playground/validator, respuestas JSON consistentes con formato unificado para éxito y error, y documentación inline de cada endpoint.

6. **Implementación de casos de uso**: Se implementaron 5 módulos de aplicación, cada uno con sus DTOs de entrada/salida, validación de reglas de negocio, manejo de errores con errores centinela, y tests unitarios con el patrón table-driven de Go cubriendo escenarios de éxito y error.

7. **Implementación del sistema de temas**: Se crearon archivos CSS con más de 30 variables personalizadas cubriendo colores, tipografías, espaciados, radios de borde y sombras para los modos claro y oscuro. El ThemeProvider de React detecta automáticamente la preferencia del sistema operativo y persiste la selección manual del usuario en localStorage.

8. **Desarrollo de componentes plantilla**: Se implementaron componentes UI reutilizables siguiendo el principio de presentación pura: no realizan llamadas a APIs, no gestionan estado global, y reciben toda la información que necesitan mediante props. Cada componente incluye múltiples variantes visuales y soporte completo para el sistema de temas.

9. **Implementación del cliente API con renovación automática de tokens**: Se configuró axios con interceptores de petición (inyección del token JWT) y respuesta (manejo de 401 con renovación automática). Se implementó una cola de peticiones fallidas durante la renovación para evitar pérdida de solicitudes y proporcionar una experiencia de usuario fluida.

### Infraestructura

10. **Dockerización multi-stage**: Se crearon Dockerfiles optimizados para cada servicio. El backend utiliza compilación en dos etapas (golang:1.23-alpine → distroless) reduciendo la imagen final a ~15 MB. El frontend utiliza build con Node 22 y sirve los archivos estáticos con Nginx, optimizando la caché y la compresión.

11. **Orquestación con Docker Compose**: Se configuraron tres servicios (backend, frontend, db) con redes separadas (interna para la comunicación backend-DB, web para frontend), health checks con tiempos de espera configurables, volúmenes persistentes para la base de datos, y variables de entorno externalizadas en archivo `.env`.

### Testing

12. **Pruebas de integración con testcontainers-go**: Se implementaron tests que levantan contenedores PostgreSQL reales desde el código de prueba, ejecutan las migraciones automáticamente, y verifican el comportamiento de los repositorios contra una base de datos auténtica, garantizando que las queries SQL generadas por SQLC funcionan correctamente con PostgreSQL 16.

13. **Pruebas unitarias con table-driven tests**: Cada caso de uso se probó con el patrón table-driven de Go, que permite definir múltiples escenarios (éxito, error, casos límite) en una estructura de datos y ejecutarlos todos con un único bucle de test. Esto mejora la cobertura, la legibilidad y el mantenimiento de los tests.

---

## Resultados del Proyecto

### Backend

1. **Módulo de autenticación funcional**: Registro de nuevos usuarios, inicio de sesión con verificación de credenciales, renovación automática de tokens JWT, y cierre de sesión. Los tokens de acceso expiran en 15 minutos y los de refresco en 7 días. El módulo incluye verificación de cuenta activa y diferenciación de roles (super_admin, user).

2. **Módulo de gestión de deudas**: CRUD completo de deudas con campos para monto total, monto restante, tasa de interés, fechas de inicio y fin, y estado (activa, pagada, liquidada). Seguimiento mes a mes con registro de monto debido, monto pagado, estado de pago y notas. Resumen anual por deuda con totales pagados y pendientes.

3. **Módulo de gastos recurrentes**: CRUD de suscripciones y gastos fijos con soporte para frecuencias (mensual, anual, semanal, quincenal), día de vencimiento, y estado (activo, cancelado). Seguimiento mensual de pagos con registro de fecha y monto. Proyección de gastos futuros basada en gastos activos.

4. **Módulo de contabilidad**: Registro de transacciones de ingresos y egresos con categorización, cálculo de balance mensual (total ingresos, total egresos, balance neto, pago de deudas), flujo de caja por periodos configurables, y proyecciones financieras basadas en datos históricos y estimaciones del usuario.

5. **Módulo de administración**: Listado paginado de usuarios con filtros, creación de nuevos usuarios, modificación de datos y permisos, activación/desactivación de cuentas. Acceso restringido exclusivamente al super-usuario mediante middleware de verificación de rol.

6. **API REST versionada**: 25+ endpoints organizados bajo `/api/v1/` con autenticación JWT obligatoria en rutas protegidas, respuestas JSON consistentes con formato `{ data, meta, error }`, validación de entrada en todos los endpoints, y errores descriptivos sin exponer detalles internos.

### Frontend

7. **Sistema de temas claro/oscuro**: Implementado con más de 30 variables CSS personalizadas y ThemeProvider de React. Detecta automáticamente la preferencia del sistema operativo mediante `prefers-color-scheme`, permite cambio manual con un toggle, y persiste la selección en localStorage.

8. **15 componentes plantilla reutilizables**: Button (4 variantes, 3 tamaños), Card (3 variantes), Table (3 variantes con ordenación), Modal (4 tamaños), Input (2 variantes con label, error, helper), Select (2 variantes), Badge (4 variantes), Loading (3 variantes), EmptyState, StatsCard (3 tendencias), PageHeader, Tabs (3 variantes), Switch, DataTable (con paginación cliente/servidor). Todos con soporte completo de temas mediante variables CSS.

9. **Cliente API con renovación automática de tokens**: Implementado con axios interceptors que inyectan automáticamente el token JWT en cada petición, detectan respuestas 401, renuevan el token usando el refresh token almacenado, reencolan las peticiones fallidas, y proporcionan una experiencia de usuario sin interrupciones.

10. **Contexto de autenticación completo**: AuthContext de React que expone usuario actual, estado de carga, funciones de login/register/logout, y banderas de autenticación y rol de administrador. Integrado con el sistema de rutas protegidas del router.

11. **Seis páginas funcionales**: Dashboard con cards de resumen financiero y gráfica de balance de últimos 6 meses; Deudas con tabla CRUD y vista de seguimiento mensual; Gastos Recurrentes con lista y toggle de pago por mes; Contabilidad con registro de transacciones y balance mensual; Planificación con proyecciones y estimaciones; Admin con tabla de usuarios y modal de edición.

### Base de Datos

12. **Esquema PostgreSQL de 12 tablas**: Migraciones versionadas con up y down, índices en todas las claves foráneas y columnas de búsqueda frecuente, constraints CHECK para valores enumerados, y claves primarias UUID para escalabilidad horizontal.

13. **Cifrado de datos sensibles**: Los campos email, full_name y phone se almacenan cifrados con AES-256-GCM con nonce aleatorio de 12 bytes. Las contraseñas se almacenan con bcrypt con costo 12. El cifrado y descifrado ocurre exclusivamente en la capa de repositorio, transparente para los casos de uso.

### Infraestructura

14. **Tres contenedores Docker**: Backend en imagen distroless de ~15 MB, frontend servido con Nginx con compresión gzip y caché de archivos estáticos, y PostgreSQL 16 con volumen persistente y health checks.

15. **Docker Compose funcional**: Orquestación completa con definición de servicios, redes separadas (interna para backend-DB, web para frontend), dependencias con condition `service_healthy`, variables de entorno externalizadas, y perfiles para desarrollo y producción.

### Documentación

16. **Documentación técnica completa**: Prompt de desarrollo con 12 fases detalladas que cubren desde la preparación del entorno hasta el despliegue en producción, incluyendo especificación de arquitectura, diseño de base de datos, guía de estilo de código, y configuración de herramientas.

---

## Conclusiones y Recomendaciones

### Conclusiones

1. **La arquitectura hexagonal demostró ser efectiva** para lograr el desacoplamiento entre funcionalidades. Cada módulo pudo desarrollarse de forma independiente, y las pruebas unitarias de los casos de uso se realizaron sin necesidad de infraestructura real, simplemente mockeando las interfaces de repositorio. Esto validó que la inversión inicial en definir interfaces claras se amortiza rápidamente con la facilidad de testing y mantenimiento.

2. **SQLC combinado con pgx ofrece una experiencia de desarrollo superior** en comparación con ORMs tradicionales. El código generado es completamente tipado, eficiente en tiempo de ejecución, y seguro contra inyección SQL por construcción. La principal ventaja es que el desarrollador escribe SQL directamente, manteniendo control total sobre las queries, mientras que SQLC elimina el trabajo repetitivo de mapear resultados a structs Go.

3. **El sistema de componentes plantilla con variables CSS cumple su objetivo** de permitir cambios visuales globales. Modificar el tema claro/oscuro o ajustar la tipografía requiere únicamente actualizar las variables CSS en dos archivos, sin tocar ningún componente individual. Esto reduce drásticamente el esfuerzo de mantenimiento visual y permite experimentar con diferentes identidades visuales sin costos de desarrollo.

4. **La combinación de JWT con renovación automática proporciona una experiencia de usuario fluida** sin sacrificar seguridad. Los tokens de corta duración (15 minutos) limitan la ventana de exposición ante posibles fugas, mientras que la renovación automática mediada por interceptores de axios ocurre de forma transparente para el usuario, que permanece autenticado durante toda su sesión sin intervención manual.

5. **El cifrado a nivel de aplicación es la estrategia correcta** para este proyecto. Permite utilizar funciones estándar de PostgreSQL sin depender de extensiones como pgcrypto, mantiene la lógica de cifrado cerca del código que la necesita, y facilita la rotación de claves. La desventaja es que las columnas cifradas no pueden ser indexadas ni utilizadas en búsquedas SQL directas, pero para el volumen de datos de un usuario individual esto no representa un problema práctico.

6. **Docker multi-stage reduce significativamente el tamaño de las imágenes** de producción. La imagen del backend pasa de ~1.2 GB (imagen de compilación) a ~15 MB (imagen distroless), mejorando los tiempos de descarga, despliegue y reduciendo la superficie de ataque. La imagen de producción contiene únicamente el binario compilado, sin shell, sin compilador y sin herramientas del sistema.

### Recomendaciones

1. **Expansión a aplicación móvil**: La API REST y los componentes plantilla del frontend están diseñados con la reutilización multiplataforma en mente. Se recomienda desarrollar una app móvil con React Native que consuma la misma API backend y adapte visualmente los componentes existentes al ecosistema móvil, maximizando el código compartido.

2. **Implementación de caché con Redis**: Los balances mensuales, los resúmenes anuales y las proyecciones implican agregaciones de múltiples tablas que pueden ser computacionalmente intensivas a medida que crece el volumen de datos. Se recomienda implementar Redis como capa de caché en los endpoints de consulta más pesados, con invalidación controlada mediante eventos de escritura.

3. **Automatización de backups**: Implementar un sistema de backup automático de PostgreSQL con `pg_dump`, programado mediante cron o systemd timer, con rotación de backups y almacenamiento en almacenamiento externo (S3, backup remoto). Incluir verificación periódica de integridad de los backups mediante restauración en un entorno aislado.

4. **Monitoreo y logging estructurado**: Incorporar logging estructurado con `slog` (log/slog de la stdlib de Go 1.21+) en el backend, exportación de métricas a Prometheus, y dashboards en Grafana para monitoreo en producción. Esto permite detectar problemas de rendimiento, cuellos de botella en la base de datos, y patrones de uso anómalos.

5. **Ampliación del sistema de roles**: Actualmente solo existen dos roles (super_admin y user). Para escenarios donde múltiples personas gestionen las finanzas de un hogar, se recomienda implementar roles adicionales como admin (gestión sin acceso a datos cifrados de otros usuarios), viewer (solo lectura), y editor (lectura y escritura limitada).

6. **Soporte multi-moneda**: Añadir soporte para múltiples monedas con tasas de conversión configurables manualmente o mediante integración con APIs de tipos de cambio. Almacenar cada transacción en su moneda original y convertir a la moneda base del usuario para los cálculos de balance y proyecciones.

7. **Importación de datos bancarios**: Implementar importación de extractos bancarios en formatos CSV y OFX, con matching automático de transacciones contra categorías existentes, detección de duplicados, y procesamiento por lotes. Esto reduciría significativamente la fricción inicial de registro de datos históricos.

8. **Versión offline-first para móvil**: Para la versión móvil, considerar una arquitectura offline-first con SQLite como almacenamiento local, sincronización diferida con el servidor, y manejo de conflictos mediante resolución por última escritura o intervención del usuario. Esto permitiría usar la app sin conexión a internet en zonas sin cobertura.

---

## Recursos y Medios Necesarios

### Recursos de Software

| Recurso | Versión | Propósito |
|---------|---------|-----------|
| Go | 1.23+ | Lenguaje de programación del backend |
| Gin | v1.10+ | Framework HTTP para Go |
| pgx | v5 | Driver nativo de PostgreSQL para Go |
| SQLC | latest | Generación de código Go desde SQL |
| golang-migrate | v4 | Sistema de migraciones de base de datos |
| golang-jwt | v5 | Creación y validación de tokens JWT |
| golang.org/x/crypto | latest | Hashing bcrypt para contraseñas |
| go-playground/validator | v10 | Validación de estructuras y requests |
| testcontainers-go | latest | Tests de integración con contenedores |
| React | 19 | Librería de interfaz de usuario |
| TypeScript | 5.x | Superset tipado de JavaScript |
| Vite | 6 | Herramienta de build y dev server |
| TailwindCSS | 4 | Framework de estilos utilitarios |
| axios | latest | Cliente HTTP con interceptors |
| react-router-dom | v7 | Routing para SPA |
| recharts | latest | Librería de gráficas |
| PostgreSQL | 16 | Sistema de gestión de bases de datos |
| Docker | latest | Plataforma de contenedores |
| Docker Compose | v2 | Orquestación multi-contenedor |
| Vitest | latest | Framework de testing para frontend |
| Testing Library | latest | Testing de componentes React |
| Nginx | stable | Servidor web para frontend en producción |
| Caddy (opcional) | 2 | Proxy inverso con SSL automático |

### Recursos de Hardware

| Recurso | Especificación | Propósito |
|---------|---------------|-----------|
| Servidor de desarrollo | 4 GB RAM, 2 vCPUs, 20 GB SSD | Entorno de desarrollo y pruebas locales |
| Servidor de producción | 2 GB RAM, 1 vCPU, 10 GB SSD | Entorno de producción para usuario único o familiar |
| Servidor de producción | 4 GB RAM, 2 vCPUs, 20 GB SSD | Entorno de producción multi-usuario (10-50 cuentas) |
| Almacenamiento | SSD, mínimo 10 GB | Volumen persistente Docker para base de datos |

### Referencias Bibliográficas

1. Pike, R. (2012). *The Go Programming Language*. Addison-Wesley Professional.

2. Chen, L. (2018). *Hexagonal Architecture: Principles and Implementation*. Journal of Software Engineering, 12(3), 145-162.

3. Banks, A., & Porcello, E. (2020). *Learning React: Functional Web Development with React and Redux* (2nd ed.). O'Reilly Media.

4. PostgreSQL Global Development Group. (2024). *PostgreSQL 16 Documentation*. Recuperado de https://www.postgresql.org/docs/16/

5. Docker Inc. (2024). *Docker Documentation*. Recuperado de https://docs.docker.com/

6. Gin Contributors. (2024). *Gin Web Framework Documentation*. Recuperado de https://gin-gonic.com/docs/

7. Kelsey, T. (2017). *Cloud Native Development with Docker*. O'Reilly Media.

8. Fowler, M. (2002). *Patterns of Enterprise Application Architecture*. Addison-Wesley Professional.

9. Freeman, S., & Pryce, N. (2009). *Growing Object-Oriented Software, Guided by Tests*. Addison-Wesley Professional.

10. Evans, E. (2003). *Domain-Driven Design: Tackling Complexity in the Heart of Software*. Addison-Wesley Professional.

11. SQLC Authors. (2024). *SQLC Documentation*. Recuperado de https://docs.sqlc.dev/

12. Tailwind Labs. (2024). *TailwindCSS Documentation*. Recuperado de https://tailwindcss.com/docs

13. Golang JWT Authors. (2024). *golang-jwt Documentation*. Recuperado de https://github.com/golang-jwt/jwt

14. Testcontainers Authors. (2024). *Testcontainers for Go*. Recuperado de https://golang.testcontainers.org/

15. National Institute of Standards and Technology. (2001). *Advanced Encryption Standard (AES)*. Federal Information Processing Standards Publication 197.

16. Project JWT Contributors. (2024). *JSON Web Token RFC 7519*. Recuperado de https://jwt.io/
