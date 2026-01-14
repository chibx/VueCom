CREATE DATABASE vuecom_catalog;
CREATE DATABASE vuecom_orders;
CREATE DATABASE vuecom_users;
CREATE DATABASE vuecom_inventory;
CREATE DATABASE vuecom_payments;
CREATE DATABASE vuecom_notifications;

CREATE USER vuecom;
GRANT ALL PRIVILEGES ON DATABASE vuecom_catalog TO vuecom;
GRANT ALL PRIVILEGES ON DATABASE vuecom_orders TO vuecom;
GRANT ALL PRIVILEGES ON DATABASE vuecom_users TO vuecom;
GRANT ALL PRIVILEGES ON DATABASE vuecom_inventory TO vuecom;
GRANT ALL PRIVILEGES ON DATABASE vuecom_payments TO vuecom;
GRANT ALL PRIVILEGES ON DATABASE vuecom_notifications TO vuecom;