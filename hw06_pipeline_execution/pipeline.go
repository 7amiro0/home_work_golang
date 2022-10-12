package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func writeReadyValue(in In, done In, out Bi) {
	for value := range in {
		select {
		case <-done:
			return
		default:
		}
		select {
		case <-done:
			return
		case out <- value:
		}
	}
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)
	go func() {
		defer close(out)

		for _, task := range stages {
			select {
			case <-done:
				return
			default:
				in = task(in)
			}
		}

		select {
		case <-done:
			return
		case value, _ := <-in:
			if value != nil {
				out <- value
			}
		}

		writeReadyValue(in, done, out)
	}()
	return out
}
