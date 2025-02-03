package closer

import (
	"crypto_scam/internal/logger"
	"os"
	"os/signal"
	"sync"
)

// globalCloser — глобальный экземпляр Closer для управления graceful shutdown.
var globalCloser = New()

// Add добавляет одну или несколько функций в глобальный Closer.
// Эти функции будут выполнены при завершении работы программы.
func Add(f ...func() error) {
	globalCloser.Add(f...)
}

// Wait блокирует выполнение до тех пор,
// пока все зарегистрированные функции не будут выполнены.
func Wait() {
	globalCloser.Wait()
}

// CloseAll вызывает выполнение всех зарегистрированных функций
// и завершает работу Closer.
func CloseAll() {
	globalCloser.CloseAll()
}

// Closer — структура для управления graceful shutdown.
// Она позволяет регистрировать функции, которые будут выполнены при завершении работы.
//
// mu - мьютекс для обеспечения потокобезопасности.
// once - гарантирует однократное выполнение CloseAll.
// done - канал для сигнализации о завершении всех функций.
// funcs - слайс функций, которые будут выполнены при завершении.
type Closer struct {
	mu    sync.Mutex     //
	once  sync.Once      //
	done  chan struct{}  //
	funcs []func() error //
}

// New создает новый экземпляр Closer.
// Если переданы сигналы ОС (например, os.Interrupt), Closer автоматически
// вызовет CloseAll при получении одного из этих сигналов.
func New(sig ...os.Signal) *Closer {
	c := &Closer{done: make(chan struct{})}
	if len(sig) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, sig...) // Регистрируем сигналы ОС.
			<-ch                      // Ожидаем получения сигнала.
			signal.Stop(ch)           // Прекращаем отслеживание сигналов.
			c.CloseAll()              // Вызываем завершение.
		}()
	}
	return c
}

// Add добавляет одну или несколько функций в Closer.
// Эти функции будут выполнены при вызове CloseAll.
func (c *Closer) Add(f ...func() error) {
	c.mu.Lock()
	c.funcs = append(c.funcs, f...) // Добавляем функции в слайс.
	c.mu.Unlock()
}

// Wait блокирует выполнение горутины, в которой он был вызван,
// до тех пор, пока все зарегистрированные функции не будут выполнены
// и канал done не будет закрыт.
// Это позволяет дождаться завершения работы всех функций,
// зарегистрированных в Closer.
func (c *Closer) Wait() {
	<-c.done
}

// CloseAll выполняет все зарегистрированные функции и завершает работу Closer.
// Гарантирует, что функции будут выполнены только один раз.
func (c *Closer) CloseAll() {
	c.once.Do(func() {
		defer close(c.done) // Закрываем канал done после выполнения всех функций.

		c.mu.Lock()
		funcs := c.funcs // Копируем слайс функций.
		c.funcs = nil    // Очищаем слайс, чтобы избежать повторного выполнения.
		c.mu.Unlock()

		// Выполняем все функции асинхронно.
		errs := make(chan error, len(funcs))
		for _, f := range funcs {
			go func(f func() error) {
				errs <- f() // Выполняем функцию и отправляем ошибку в канал.
			}(f)
		}

		// Обрабатываем ошибки, если они возникли.
		for i := 0; i < cap(errs); i++ {
			if err := <-errs; err != nil {
				logger.Info("error returned from Closer:", err)
			}
		}
	})
}
