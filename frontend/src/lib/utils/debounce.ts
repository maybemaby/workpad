/**
 * Creates a debounced version of a function that delays its execution
 * @param callback The function to debounce
 * @param duration The delay in milliseconds
 * @returns A debounced function with inferred parameter and return types
 */
export function debounce<Args extends unknown[], Return>(
	callback: (...args: Args) => Return,
	duration: number
): (...args: Args) => void {
	let timeoutId: ReturnType<typeof setTimeout> | null = null;

	return function debounced(...args: Args): void {
		if (timeoutId !== null) {
			clearTimeout(timeoutId);
		}

		timeoutId = setTimeout(() => {
			callback(...args);
			timeoutId = null;
		}, duration);
	};
}
