.PHONY: update-docs check-docs

# Обновление sourcecraft-документации на основе текущего кода
# Запускает AI-агента для анализа изменений и обновления файлов
update-docs:
	@echo "=== Обновление sourcecraft документации ==="
	@echo ""
	@echo "Запустите AI-агента с задачей:"
	@echo ""
	@echo "  Проанализируй код сервера (server/) и клиента (client/)"
	@echo "  и обнови файлы:"
	@echo "    - sourcecraft/rules.md"
	@echo "    - sourcecraft/arch.md"
	@echo "    - sourcecraft/requirements.md"
	@echo ""
	@echo "=== Сравнение с последним коммитом ==="
	git diff --name-only HEAD~1..HEAD 2>/dev/null || echo "(нет предыдущего коммита для сравнения)"
	@echo ""
	@echo "=== Изменения в server/ ==="
	git diff --stat HEAD~1..HEAD -- server/ 2>/dev/null || echo "(нет)"
	@echo ""
	@echo "=== Изменения в client/ ==="
	git diff --stat HEAD~1..HEAD -- client/ 2>/dev/null || echo "(нет)"

# Проверка, нужно ли обновлять документацию
check-docs:
	@echo "=== Проверка необходимости обновления документации ==="
	@CHANGED=$$(git diff --name-only HEAD~1..HEAD -- server/ client/ 2>/dev/null); \
	if [ -n "$$CHANGED" ]; then \
		echo "Изменены файлы:"; \
		echo "$$CHANGED"; \
		echo ""; \
		echo "⚠️  Рекомендуется запустить: make update-docs"; \
	else \
		echo "✅ Изменений в коде нет, документация актуальна."; \
	fi
