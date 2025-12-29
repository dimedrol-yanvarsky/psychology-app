# Структура клиентской части

Клиент — это Create React App проект в каталоге `client/`.

## Верхний уровень

- `package.json` / `package-lock.json` — зависимости и скрипты.
- `public/` — статические файлы, `index.html`, `robots.txt`.
- `src/` — исходный код приложения.

## Краткое дерево

```
client/
  public/
  src/
    app/
      App.js
      styles/
    pages/
      dashboard/
      login/
      recommendations/
      registration/
      reviews/
      terminal/
      tests/
      tree/
    shared/
      assets/
      config/
      lib/
      model/
      ui/
    widgets/
      header/
      terminal/
    index.js
```

## src/

- `index.js` — точка входа и монтирование React-приложения.
- `app/` — оболочка приложения и глобальные стили.
  - `App.js` — корневой компонент.
  - `styles/` — глобальные CSS-стили.
- `pages/` — страницы, сгруппированные по функциональности.
  - `dashboard/`, `login/`, `recommendations/`, `registration/`, `reviews/`,
    `terminal/`, `tests/`, `tree/`.
  - Типичные подпапки: `ui/` (компоненты страницы), `api/` (HTTP-вызовы),
    `model/` (хуки/состояние), `lib/` (вспомогательные функции), `index.js` (экспорт).
- `widgets/` — крупные переиспользуемые блоки интерфейса (`header/`, `terminal/`).
- `shared/` — общие для всего приложения сущности.
  - `assets/` — изображения и шрифты.
  - `config/` — общие настройки (например, `api.js`).
  - `lib/` — утилиты и хуки (например, `useLockBodyScroll.js`).
  - `model/` — общие модели/состояние (например, `profile.js`).
  - `ui/` — переиспользуемые UI-компоненты (кнопка, модалка, алерт).
