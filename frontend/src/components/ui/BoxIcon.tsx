import type { ComponentProps, ReactElement } from 'react';

export interface BoxIconProps extends ComponentProps<'span'> {
    className?: string;
}

type IconComponent = (props: BoxIconProps) => ReactElement;

function createBoxIcon(iconClass: string): IconComponent {
    return function BoxIcon({ className = '', style, ...props }: BoxIconProps) {
        return (
            <span
                className={`inline-flex items-center justify-center ${className}`}
                style={style}
                {...props}
            >
                <i className={`bx ${iconClass} bx-icon`} />
            </span>
        );
    };
}

// Navigation & UI
export const Target = createBoxIcon('bxs-bullseye');
export const Crosshair = createBoxIcon('bx-crosshair');
export const Layers = createBoxIcon('bxs-layer');
export const Crown = createBoxIcon('bxs-crown');
export const User = createBoxIcon('bxs-user');
export const Settings = createBoxIcon('bxs-cog');
export const LogOut = createBoxIcon('bxs-log-out');
export const LogIn = createBoxIcon('bxs-log-in');
export const Palette = createBoxIcon('bxs-palette');
export const Moon = createBoxIcon('bxs-moon');
export const Sun = createBoxIcon('bxs-sun');
export const Globe = createBoxIcon('bxs-globe');
export const Home = createBoxIcon('bxs-home');
export const Bell = createBoxIcon('bxs-bell');
export const Eye = createBoxIcon('bxs-show');

// Game
export const Trophy = createBoxIcon('bxs-trophy');
export const Clock = createBoxIcon('bxs-time-five');
export const Zap = createBoxIcon('bxs-bolt');
export const InfinityIcon = createBoxIcon('bx-infinite');
export const Play = createBoxIcon('bxs-right-arrow');
export const Flame = createBoxIcon('bxs-flame');

// Actions
export const Send = createBoxIcon('bxs-send');
export const Save = createBoxIcon('bxs-save');
export const Edit2 = createBoxIcon('bxs-edit-alt');
export const RotateCcw = createBoxIcon('bx-rotate-left');
export const RefreshCw = createBoxIcon('bx-refresh');
export const X = createBoxIcon('bx-x');

// Chevrons
export const ChevronLeft = createBoxIcon('bx-chevron-left');
export const ChevronRight = createBoxIcon('bx-chevron-right');
export const ChevronUp = createBoxIcon('bx-chevron-up');
export const ChevronDown = createBoxIcon('bx-chevron-down');

// Alerts
export const AlertCircle = createBoxIcon('bxs-error-circle');
export const AlertTriangle = createBoxIcon('bxs-error');
export const CheckCircle = createBoxIcon('bxs-check-circle');
export const Info = createBoxIcon('bxs-info-circle');

// Stats
export const TrendingUp = createBoxIcon('bx-trending-up');

// Type for icon components (replaces LucideIcon)
export type BoxIconType = IconComponent;