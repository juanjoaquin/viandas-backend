# Sistema de Gestión de Viandas

## Objetivo

Digitalizar el proceso actual de gestión de viandas.

La empresa trabaja con:

* Menús semanales.
* Clientes (personas o empresas).
* Deliverys.
* Producción diaria.
* Extras (ensaladas y sándwiches).

No se contempla en esta versión:

* Facturación.
* Pagos.
* Stock.
* Impresión.

---

# Roles

## ADMIN

Responsabilidades:

* Gestionar usuarios.
* Gestionar clientes.
* Gestionar deliverys.
* Gestionar platos.
* Gestionar extras.
* Configurar menús semanales.
* Consultar producción.

---

## EMPLOYEE

Responsabilidades:

* Consultar menú semanal.
* Registrar producción diaria.
* Modificar producción diaria.
* Consultar clientes.
* Consultar deliverys.

No puede:

* Gestionar usuarios.
* Gestionar platos.
* Gestionar extras.
* Configurar menús semanales.

---

# Modelo de Datos

## users

Representa usuarios del sistema.

```sql
id UUID PRIMARY KEY

name VARCHAR(100)

email VARCHAR(255) UNIQUE

password_hash TEXT

role VARCHAR(20)

active BOOLEAN

created_at TIMESTAMP
updated_at TIMESTAMP
```

### Roles

```text
ADMIN
EMPLOYEE
```

---

## customers

Representa empresas o personas que realizan pedidos.

```sql
id UUID PRIMARY KEY

name VARCHAR(255)

type VARCHAR(20)

phone VARCHAR(100)

address TEXT

created_at TIMESTAMP
updated_at TIMESTAMP
```

### Tipos

```text
COMPANY
PERSON
```

---

## deliveries

Representa repartidores.

```sql
id UUID PRIMARY KEY

name VARCHAR(255)

active BOOLEAN

created_at TIMESTAMP
updated_at TIMESTAMP
```

---

## dishes

Catálogo de platos.

Un plato pertenece siempre a un único tipo de menú.

```sql
id UUID PRIMARY KEY

name VARCHAR(255)

description TEXT

menu_type VARCHAR(30)

active BOOLEAN

created_at TIMESTAMP
updated_at TIMESTAMP
```

### Menu Types

```text
TRADITIONAL
HEALTHY
VEGETARIAN
```

### Ejemplos

```text
Milanesa con puré
TRADITIONAL

Canelones de carne
TRADITIONAL

Tarta de jamón y queso
HEALTHY

Risotto de hongos
VEGETARIAN
```

---

## extra_products

Productos disponibles todos los días.

```sql
id UUID PRIMARY KEY

name VARCHAR(255)

category VARCHAR(30)

active BOOLEAN

created_at TIMESTAMP
updated_at TIMESTAMP
```

### Categorías

```text
SALAD
SANDWICH
```

### Ejemplos

```text
Ensalada César

Ensalada Completa

Sándwich de Pollo

Sándwich Jamón y Queso
```

---

# Menú Semanal

El administrador configura semanalmente qué plato corresponde a cada menú.

---

## week_menus

Representa una semana.

```sql
id UUID PRIMARY KEY

week_start_date DATE

created_by UUID REFERENCES users(id)

created_at TIMESTAMP
updated_at TIMESTAMP
```

### Ejemplo

```text
Semana del 15/06/2026
```

---

## week_menu_items

Representa la configuración de un día.

```sql
id UUID PRIMARY KEY

week_menu_id UUID REFERENCES week_menus(id)

menu_date DATE

traditional_dish_id UUID REFERENCES dishes(id)

healthy_dish_id UUID REFERENCES dishes(id)

vegetarian_dish_id UUID REFERENCES dishes(id)

created_at TIMESTAMP
updated_at TIMESTAMP
```

### Ejemplo

```text
Lunes

Tradicional:
Milanesa con puré

Saludable:
Tarta de jamón y queso

Vegetariano:
Risotto de hongos
```

---

# Producción Diaria

Representa la planilla principal utilizada por los empleados.

Importante:

Los empleados NO seleccionan platos.

Los empleados únicamente registran cantidades por tipo de menú.

El plato asociado se obtiene automáticamente desde el menú semanal correspondiente a la fecha.

---

## daily_productions

Representa una fila de la planilla.

```sql
id UUID PRIMARY KEY

production_date DATE

customer_id UUID REFERENCES customers(id)

delivery_id UUID REFERENCES deliveries(id)

traditional_qty INTEGER DEFAULT 0

healthy_qty INTEGER DEFAULT 0

vegetarian_qty INTEGER DEFAULT 0

notes TEXT

created_by UUID REFERENCES users(id)

created_at TIMESTAMP
updated_at TIMESTAMP
```

### Ejemplo

```text
Empresa A

Tradicional: 8
Saludable: 2
Vegetariano: 0

Delivery: Juan
```

---

## daily_production_extras

Representa extras asociados a un registro de producción.

```sql
id UUID PRIMARY KEY

daily_production_id UUID REFERENCES daily_productions(id)

extra_product_id UUID REFERENCES extra_products(id)

quantity INTEGER

created_at TIMESTAMP
updated_at TIMESTAMP
```

### Ejemplo

```text
Empresa A

Ensalada César x2

Sándwich de Pollo x1
```

---

# Relaciones

```text
User
├── WeekMenu
└── DailyProduction

Customer
└── DailyProduction

Delivery
└── DailyProduction

WeekMenu
└── WeekMenuItem

Dish
└── WeekMenuItem

DailyProduction
└── DailyProductionExtra

ExtraProduct
└── DailyProductionExtra
```

---

# Reglas de Negocio

## Menús

Existen únicamente tres tipos de menú:

```text
TRADITIONAL
HEALTHY
VEGETARIAN
```

No existe CRUD para tipos de menú.

Son valores fijos del sistema.

---

## Platos

Un plato siempre pertenece al mismo menú.

Ejemplo:

```text
Milanesa con puré
→ TRADITIONAL

Tarta de jamón y queso
→ HEALTHY
```

No puede cambiar de categoría.

---

## Menú Semanal

Cada semana debe tener:

* Lunes
* Martes
* Miércoles
* Jueves
* Viernes

Cada día debe tener:

* 1 plato tradicional
* 1 plato saludable
* 1 plato vegetariano

---

## Producción

Cada registro representa un cliente para una fecha determinada.

Ejemplo:

```text
Empresa A

Tradicional: 8
Saludable: 2

Delivery: Juan
```

---

## Extras

Los extras deben almacenarse individualmente para poder obtener estadísticas y totales.

Ejemplo:

```text
Ensalada César x5

Ensalada Completa x2

Sándwich de Pollo x4
```

---

# Consultas Principales

## Menú del día

Entrada:

```text
15/06/2026
```

Salida:

```text
Tradicional
Milanesa con puré

Saludable
Tarta de jamón y queso

Vegetariano
Risotto de hongos
```

---

## Producción del día

```text
Empresa A

Tradicional: 8
Saludable: 2

Delivery: Juan
```

```text
Empresa B

Tradicional: 4

Delivery: Pedro
```

---

## Totales de cocina

```text
Tradicional: 12

Saludable: 2

Vegetariano: 0
```

---

## Totales de extras

```text
Ensalada César: 4

Ensalada Completa: 2

Sándwich de Pollo: 3
```
