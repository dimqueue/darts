/** @type {import('tailwindcss').Config} */
export default {
    content: [
        "./index.html",
        "./src/**/*.{js,ts,jsx,tsx}",
    ],
    theme: {
        extend: {
            colors: {
                primary: {
                    purple: '#A855F7',
                    blue: '#3B82F6',
                    green: '#10B981',
                },
                secondary: {
                    purple: '#EC4899',
                    blue: '#06B6D4',
                    green: '#3B82F6',
                },
            },
            backgroundImage: {
                'gradient-purple': 'linear-gradient(135deg, #A855F7 0%, #EC4899 100%)',
                'gradient-blue': 'linear-gradient(135deg, #3B82F6 0%, #06B6D4 100%)',
                'gradient-green': 'linear-gradient(135deg, #10B981 0%, #3B82F6 100%)',
            },
            borderRadius: {
                'card': '24px',
            },
            boxShadow: {
                'card': '0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03)',
            },
        },
    },
    plugins: [],
}