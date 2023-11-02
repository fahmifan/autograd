export type Result<T, E = string> =
	| { value: T; ok: true }
	| { error: E; ok: false };

export type ResultPromise<T, E = string> = Promise<Result<T, E>>;

export function Err<E>(why: E): Result<never, E> {
	return {
		error: why,
		ok: false,
	};
}

export function OK<T>(value: T): Result<T, never> {
	return {
		value: value,
		ok: true,
	};
}

export async function ResultFromPromise<T, E = string>(
	promise: Promise<T>,
): ResultPromise<T, E> {
	return promise.then(OK).catch(Err);
}
