FROM php:8.4-cli-bookworm AS base
RUN --mount=type=secret,id=custom_ca \
    if [ -f /run/secrets/custom_ca ]; then cp /run/secrets/custom_ca /usr/local/share/ca-certificates/environment.crt; update-ca-certificates; fi \
    && sed -i 's|http://deb.debian.org|https://deb.debian.org|g' /etc/apt/sources.list.d/debian.sources \
    && apt-get update && apt-get install -y --no-install-recommends git unzip libicu-dev libpq-dev libzip-dev libonig-dev libxml2-dev libsqlite3-dev \
    && docker-php-ext-install -j2 intl pdo_pgsql pdo_sqlite zip mbstring bcmath \
    && rm -rf /var/lib/apt/lists/*
COPY --from=composer:2 /usr/bin/composer /usr/local/bin/composer
WORKDIR /app
FROM base AS development
COPY . .
ARG COMPOSER_INSTALL_MODE=dist
RUN --mount=type=cache,id=observatorio-composer,target=/root/.composer/cache composer install --no-interaction --prefer-install="$COMPOSER_INSTALL_MODE" --no-progress \
    && find /app -type d -exec chmod 755 {} + \
    && find /app -type f -exec chmod a+r {} + \
    && chown -R www-data:www-data storage bootstrap/cache
USER www-data
EXPOSE 8000
CMD ["php", "artisan", "serve", "--host=0.0.0.0", "--port=8000"]
