<!-- ...existing code... -->

## 🎮 Game TODO: Dynamic Window Render (сессии по 30 минут)

Правила:
- 1 квест = ~30 минут.
- Закрыл квест → поставил ✅ и начислил XP.
- Каждые 100 XP = Level Up.

Текущий прогресс: **0 XP / Level 0**

### 🟢 Level 1 — База рендера (простые победы)

- [ ] **Q1: RenderContext MVP** (+20 XP)  
  **Цель:** описать минимальный контекст для рендера окна.  
  **DoD:** есть структура/контракт с полями: `dialogData`, `widgetData`, `startData`, `middlewareData`, `state`.

- [ ] **Q2: Skeleton window.Render()** (+20 XP)  
  **Цель:** сделать каркас `Render()` без полной логики.  
  **DoD:** в `Render()` есть шаги-пустышки: load data → call getter → build widgets → return render result.

- [ ] **Q3: Merge policy данных** (+20 XP)  
  **Цель:** зафиксировать порядок merge данных для шаблона.  
  **DoD:** документирован и применён порядок, например:  
  `getterData > dialogData > startData > middlewareData`.

- [ ] **Q4: Лог рендера** (+20 XP)  
  **Цель:** добавить debug-лог ключевых шагов.  
  **DoD:** логируются state, window id, source данных, количество кнопок.

- [ ] **Q5: Smoke-check сценарий** (+20 XP)  
  **Цель:** руками пройти один сценарий callback → render.  
  **DoD:** после callback окно перерисовывается без падения.

---

### 🟡 Level 2 — Динамические кнопки

- [ ] **Q6: Button model from getter** (+25 XP)  
  **Цель:** getter возвращает список кнопок (label + callback payload).  
  **DoD:** кнопки строятся из данных getter, а не захардкожены.

- [ ] **Q7: Уникальная callback_data** (+25 XP)  
  **Цель:** убрать коллизии callback_data.  
  **DoD:** добавлен формат, например: `dlg:{dialogId}:w:{widgetId}:a:{action}:i:{idx}`.

- [ ] **Q8: Widget handler mapping** (+25 XP)  
  **Цель:** сопоставить callback → handler виджета.  
  **DoD:** роутер/диспетчер находит нужный handler по callback_data.

- [ ] **Q9: Dialog.collectHandlers()** (+25 XP)  
  **Цель:** диалог отдаёт все handler’ы окон.  
  **DoD:** есть метод сбора handler’ов из window/widget и регистрация без ручного списка.

---

### 🟠 Level 3 — Интеграция с Router/DialogManager

- [ ] **Q10: Router.RegisterDialog(dialog)** (+30 XP)  
  **Цель:** регистрировать диалог одной командой.  
  **DoD:** `RegisterDialog` подтягивает handler’ы из диалога автоматически.

- [ ] **Q11: Render после start/switchTo** (+30 XP)  
  **Цель:** на `dialogManager.start()/switchTo()` всегда рендерится окно текущего state.  
  **DoD:** нет “пустого” state без окна; после переключения сразу UI.

- [ ] **Q12: Callback flow end-to-end** (+40 XP)  
  **Цель:** полный цикл `UPDATE -> RENDER -> EDIT/SEND`.  
  **DoD:** callback меняет данные, `Render()` собирает контекст, сообщение редактируется/отправляется, storage обновляется.

---

### 🏁 Boss Fight (1 час, опционально)

- [ ] **B1: Мини-демо “Выбор конфига VPN”** (+60 XP)  
  **DoD:**  
  1) список конфигов грузится динамически,  
  2) кнопки рисуются из getter,  
  3) выбор кнопки меняет selected state,  
  4) окно корректно перерисовывается.

---

## Награды

- 100 XP: ☕ coffee reward  
- 200 XP: 🎧 30 минут музыки без чувства вины  
- 300 XP: ✅ “Render Architect” unlocked