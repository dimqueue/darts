import { memo } from 'react';
import { Trophy } from '../ui/BoxIcon';
import Card from '../ui/Card';

interface WinMessageProps {
    guessCount: number;
}

export default memo(function WinMessage({ guessCount }: WinMessageProps) {
    return (
        <Card className="bg-theme-gradient text-white text-center">
            <Trophy className="w-12 h-12 mx-auto mb-2" aria-hidden="true" />
            <h2 className="text-2xl font-bold mb-1">Congratulations!</h2>
            <p className="text-lg opacity-90" role="status">
                You found the word in {guessCount} guesses!
            </p>
        </Card>
    );
});
