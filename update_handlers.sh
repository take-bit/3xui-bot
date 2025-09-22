#!/bin/bash

# Обновляем все обработчики
for file in internal/telegram/handlers/*.go; do
    if [ "$file" != "internal/telegram/handlers/start.go" ] && [ "$file" != "internal/telegram/handlers/help.go" ]; then
        # Добавляем UseCaseManager в структуру
        sed -i '' '/type.*Handler struct {/,/}/ {
            /telegram\.BaseHandlerInterface/a\
	UseCaseManager *usecase.UseCaseManager
        }' "$file"
        
        # Обновляем конструктор
        sed -i '' 's/BaseHandler: telegram\.NewBaseHandler(useCaseManager, nil),/UseCaseManager: useCaseManager,/g' "$file"
    fi
done
