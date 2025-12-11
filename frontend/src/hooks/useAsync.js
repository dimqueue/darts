import { useState, useCallback, useRef, useEffect } from 'react';

export function useAsync(asyncFunction, options = {}) {
    const { immediate = false, onSuccess, onError } = options;

    const [state, setState] = useState({
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
        async (...args) => {
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

                setState({ data: null, loading: false, error });
                onError?.(error);

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
            execute();
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


export function useMutation(mutationFn, options = {}) {
    return useAsync(mutationFn, { ...options, immediate: false });
}

export default useAsync;
