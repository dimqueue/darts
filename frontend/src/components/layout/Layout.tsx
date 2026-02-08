import type { ReactNode } from 'react';
import Navbar from './Navbar';
import DemoToast from './DemoToast';

interface LayoutProps {
    children: ReactNode;
}

export default function Layout({ children }: LayoutProps) {
    return (
        <div className="h-screen flex flex-col bg-theme-page overflow-hidden">
            <Navbar />
            <a
                href="#main-content"
                className="sr-only focus:not-sr-only focus:absolute focus:top-20 focus:left-4 focus:z-50 focus:px-4 focus:py-2 focus:bg-white focus:text-gray-900 focus:rounded-lg focus:shadow-lg"
            >
                Skip to main content
            </a>
            <div className="flex-1 overflow-y-auto overflow-x-hidden scrollbar-thin">
                <main
                    id="main-content"
                    className="container mx-auto px-4 py-6 max-w-6xl"
                >
                    {children}
                </main>
            </div>
            <DemoToast />
        </div>
    );
}
