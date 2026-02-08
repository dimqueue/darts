import type { ReactNode } from 'react';
import Navbar from './Navbar';
import DemoToast from './DemoToast';

interface LayoutProps {
    children: ReactNode;
}

export default function Layout({ children }: LayoutProps) {
    return (
        <div className="min-h-screen bg-theme-page">
            <Navbar />
            <a
                href="#main-content"
                className="sr-only focus:not-sr-only focus:absolute focus:top-20 focus:left-4 focus:z-50 focus:px-4 focus:py-2 focus:bg-white focus:text-gray-900 focus:rounded-lg focus:shadow-lg"
            >
                Skip to main content
            </a>
            <main id="main-content" className="container mx-auto px-4 py-8 max-w-6xl">
                {children}
            </main>
            <DemoToast />
        </div>
    );
}
