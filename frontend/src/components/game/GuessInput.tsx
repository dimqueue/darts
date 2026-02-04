import { forwardRef } from 'react';
import { Send } from 'lucide-react';
import Card from '../ui/Card';
import Button from '../ui/Button';
import Input from '../ui/Input';
import ErrorAlert from '../ui/ErrorAlert';

interface GuessInputProps {
    value: string;
    onChange: (value: string) => void;
    onSubmit: () => void;
    loading: boolean;
    disabled: boolean;
    error: string;
}

const GuessInput = forwardRef<HTMLInputElement, GuessInputProps>(function GuessInput(
    { value, onChange, onSubmit, loading, disabled, error },
    ref
) {
    const handleKeyPress = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key === 'Enter' && !loading) {
            onSubmit();
        }
    };

    return (
        <Card>
            <div className="flex gap-3">
                <Input
                    ref={ref}
                    value={value}
                    onChange={(e) => onChange(e.target.value)}
                    onKeyPress={handleKeyPress}
                    placeholder="Enter your guess..."
                    disabled={loading || disabled}
                    className="flex-1"
                />
                <Button
                    onClick={onSubmit}
                    disabled={!value.trim() || disabled}
                    loading={loading}
                    icon={Send}
                >
                    Send
                </Button>
            </div>

            {error && <ErrorAlert message={error} className="mt-3" />}
        </Card>
    );
});

export default GuessInput;
