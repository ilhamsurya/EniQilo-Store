CREATE TABLE "users" (
  "user_id" SERIAL PRIMARY KEY,
  "email" varchar unique not null,
  "name" varchar not null,
  "phone_number" varchar unique not null,
  "password" varchar not null,
  "salt" varchar not null,
  "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "products" (
  "id_product" SERIAL PRIMARY KEY,
  "name" varchar NOT NULL,
  "sku" varchar NOT NULL,
  "category" varchar NOT NULL,
  "notes" varchar NOT NULL,
  "price" decimal NOT NULL,
  "stock" int NOT NULL,
  "location" varchar NOT NULL,
  "is_available" boolean NOT NULL,
  "created_at" timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "product_images" (
  "id_image" SERIAL PRIMARY KEY,
  "id_product" integer,
  "image_url" text NOT NULL,
  "created_at" timestamp DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE "product_images" ADD FOREIGN KEY ("id_product") REFERENCES "products" ("id_product");
