package bus

func DeadLetterQueue(queue string) string {
	return queue + "/$DeadLetterQueue"
}
