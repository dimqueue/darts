import { Trophy, Crosshair, InfinityIcon, Clock, Zap, Shield, Group } from '../components/ui/BoxIcon';
import type { GameMode, GameModeId, MultiplayerMode } from '../types/game';

export const GAME_MODES: Record<GameModeId, GameMode> = {
    daily: {
        id: 'daily',
        name: 'Daily Challenge',
        description: 'A new word every day. Compare with others!',
        icon: Trophy,
        available: true,
        path: '/game/daily',
        gradient: 'from-amber-400 to-orange-500',
        rated: false,
        persistent: true,
        canContinue: true,
        concurrent: true,
        hasLeaderboard: true,
    },
    competitive: {
        id: 'competitive',
        name: 'Competitive',
        description: 'Ranked games that affect your global score.',
        icon: Crosshair,
        available: true,
        path: '/game/competitive',
        gradient: 'from-red-500 to-rose-600',
        rated: true,
        persistent: true,
        canContinue: true,
        concurrent: true,
        hasLeaderboard: true,
    },
    endless: {
        id: 'endless',
        name: 'Endless',
        description: 'Play at your own pace with no pressure.',
        icon: InfinityIcon,
        available: true,
        path: '/game/endless',
        gradient: 'from-emerald-400 to-teal-500',
        rated: false,
        persistent: true,
        canContinue: true,
        concurrent: true,
        hasLeaderboard: false,
    },
    'time-attack': {
        id: 'time-attack',
        name: 'Time Attack',
        description: 'Race against the clock to find the word.',
        icon: Clock,
        available: false,
        path: '/game/time-attack',
        gradient: 'from-blue-500 to-cyan-500',
        rated: false,
        persistent: false,
        canContinue: false,
        concurrent: false,
        hasLeaderboard: false,
        comingSoon: true,
    },
    practice: {
        id: 'practice',
        name: 'Practice',
        description: 'Warm up without affecting your stats.',
        icon: Zap,
        available: true,
        path: '/game/practice',
        gradient: 'from-violet-500 to-purple-600',
        rated: false,
        persistent: false,
        canContinue: false,
        concurrent: false,
        hasLeaderboard: false,
    },
};

export const MULTIPLAYER_MODES: MultiplayerMode[] = [
    {
        id: '1v1-duel',
        name: '1v1 Duel',
        description: 'Challenge a friend head-to-head.',
        icon: Group,
        comingSoon: true,
    },
    {
        id: 'arena',
        name: 'Arena',
        description: 'Compete against multiple players live.',
        icon: Shield,
        comingSoon: true,
    },
];

const VALID_MODE_IDS = new Set<string>(Object.keys(GAME_MODES));

export function getGameMode(mode: GameModeId): GameMode {
    return GAME_MODES[mode];
}

export function isValidGameMode(mode: string | undefined): mode is GameModeId {
    return typeof mode === 'string' && VALID_MODE_IDS.has(mode);
}

export const ALL_MODE_IDS = Object.keys(GAME_MODES) as GameModeId[];

export const GAME_MODE_LIST = Object.values(GAME_MODES);

export const PLAYABLE_MODES = GAME_MODE_LIST.filter((m) => m.available && !m.comingSoon);
