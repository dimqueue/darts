import { useState, useCallback, useRef, useEffect } from 'react';

export interface AsyncState<T> {
    data: T | null;
    loading: boolean;
    error: Error | null;
}

export interface UseAsyncOptions<T> {
    immediate?: boolean;
    onSuccess?: (result: T) => void;
    onError?: (error: Error) => void;
}

export interface UseAsyncReturn<T, Args extends unknown[]> extends AsyncState<T> {
    execute: (...args: Args) => Promise<T | undefined>;
    reset: () => void;
    isIdle: boolean;
    isLoading: boolean;
    isError: boolean;
    isSuccess: boolean;
}

export function useAsync<T, Args extends unknown[] = []>(
    asyncFunction: (...args: Args) => Promise<T>,
    options: UseAsyncOptions<T> = {}
): UseAsyncReturn<T, Args> {
    const { immediate = false, onSuccess, onError } = options;

    const [state, setState] = useState<AsyncState<T>>({
        data: null,
        loading: immediate,
        error: null,
    });

    const mountedRef = useRef(true);
    const asyncFunctionRef = useRef(asyncFunction);

    useEffect(() => {
        asyncFunctionRef.current = asyncFunction;
    }, [asyncFunction]);

    useEffect(() => {
        mountedRef.current = true;
        return () => {
            mountedRef.current = false;
        };
    }, []);

    const execute = useCallback(
        async (...args: Args): Promise<T | undefined> => {
            if (!mountedRef.current) return;

            setState((prev) => ({ ...prev, loading: true, error: null }));

            try {
                const result = await asyncFunctionRef.current(...args);

                if (!mountedRef.current) return;

                setState({ data: result, loading: false, error: null });
                onSuccess?.(result);

                return result;
            } catch (error) {
                if (!mountedRef.current) return;

                const err = error instanceof Error ? error : new Error(String(error));
                setState({ data: null, loading: false, error: err });
                onError?.(err);

                throw error;
            }
        },
        [onSuccess, onError]
    );

    const reset = useCallback(() => {
        setState({ data: null, loading: false, error: null });
    }, []);

    useEffect(() => {
        if (immediate) {
            execute(...([] as unknown as Args));
        }
    }, [immediate, execute]);

    return {
        ...state,
        execute,
        reset,
        isIdle: !state.loading && !state.error && !state.data,
        isLoading: state.loading,
        isError: !!state.error,
        isSuccess: !!state.data && !state.error,
    };
}

export function useMutation<T, Args extends unknown[] = []>(
    mutationFn: (...args: Args) => Promise<T>,
    options: UseAsyncOptions<T> = {}
): UseAsyncReturn<T, Args> {
    return useAsync(mutationFn, { ...options, immediate: false });
}

export default useAsync;
