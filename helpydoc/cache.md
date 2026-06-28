## Что нужно хранить в редисе от диалога:
- chat id/user id (как ключ всего диалога)
- message id (чтобы редактировать/удалять диалог)
- state (чтобы понимать на каком состоянии диалога юзер)
- dialog_data (пользоватльские значение которые можно сериализовать)

## Актуальное уточнение
- В текущей реализации ключ dialogSession — только `user id`.
- Формат ключа в Redis: `dialog:<userID>`.
 

## Что нужно хранить в кеше от Window
    Состояния Widget только

## Что нужно хранить в кеше от Widget
### MediaMover
    CurrentPage 
    SelectedIDx

## Radio
    CheckedItemId

!!! Чтобы хранить в кеше юзерские данные виджетов ,нужно их связывать с стейтом и юзером `state:<state>userID:<userID>`


## Уровни абстракции кеширования в Widget и Dialog
### Dialog
    1) Вызов  в router dialogData и messageId

### Widget
    1) Ручной вызов кэша в ToInlineKeyboard c ключом типа `state:<state>userID:<userID>`
    2) метод getCache() (SelectedID int)/(SelectedID, CurrentPage)
    3) функция GetWidgetData (state, userID) map[string]interface{} + парсер мапы в структуру в каждый widget

### Memento(Снимок) паттрен
? Каждый виджет создает снэпшот, 


### Жизненный цикл FSM 
Приходит апдейт - из кэша узнается состояние диалога с пользователем  по user id - достается data этого состояния из кэша - создается fsm которая теперь управляет data и state диалога

тогда вопросы к storage и как его обустроить

 создать структуру когторая используя интерфейс cache внутри себя уже заточена на сохранение и оперирование в базе своим доменом 

 Будет примитивный интерфейс Storage c Get(),Set(),Delete()
 Реализация на redis и MemoryCache // MemoryCache надо сделать потокобезопасной
 
 На каждый  апдейт будет создаваться FSMСontext привязанный к диалогу, который используя примитивный Storage уже сам будет правильно парсить данные с своей структурой хранения

 С FSMContext   


Структура хранения в кэше:

UserID -> Dialog -> data
                 -> DialogSession -> CurrentState
                                  -> MessageID

       -> StateName -> data