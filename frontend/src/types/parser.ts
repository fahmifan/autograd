export function parseIntWithDefault(
	value: string | null,
	defaultValue: number,
): number {
	if (value) {
		return parseInt(value);
	}

	return defaultValue;
}
