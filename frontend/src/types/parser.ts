export function parseIntWithDefault(
	value: string | null,
	defaultValue: number,
): number {
	if (value) {
		return Number.parseInt(value);
	}

	return defaultValue;
}
