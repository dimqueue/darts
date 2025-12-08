/** @type {import('tailwindcss').Config} */
export default {
    content: [
        "./index.html",
        "./src/**/*.{js,ts,jsx,tsx}",
    ],
    darkMode: 'class',
    theme: {
        extend: {
            colors: {
                primary: {
                    purple: '#8B5CF6',
                    blue: '#3B82F6',
                    green: '#10B981',
                },
                secondary: {
                    purple: '#A78BFA',
                    blue: '#60A5FA',
                    green: '#34D399',
                },
            },
            backgroundImage: {
                'gradient-purple': 'linear-gradient(135deg, #8B5CF6 0%, #A78BFA 50%, #C4B5FD 100%)',
                'gradient-blue': 'linear-gradient(135deg, #3B82F6 0%, #60A5FA 50%, #93C5FD 100%)',
                'gradient-green': 'linear-gradient(135deg, #10B981 0%, #34D399 50%, #6EE7B7 100%)',
                'mesh-purple': 'radial-gradient(at 40% 20%, hsla(268, 82%, 85%, 0.3) 0px, transparent 50%), radial-gradient(at 80% 0%, hsla(286, 72%, 90%, 0.2) 0px, transparent 50%), radial-gradient(at 0% 50%, hsla(268, 82%, 95%, 0.3) 0px, transparent 50%)',
                'mesh-blue': 'radial-gradient(at 40% 20%, hsla(217, 91%, 85%, 0.3) 0px, transparent 50%), radial-gradient(at 80% 0%, hsla(190, 80%, 90%, 0.2) 0px, transparent 50%), radial-gradient(at 0% 50%, hsla(217, 91%, 95%, 0.3) 0px, transparent 50%)',
                'mesh-green': 'radial-gradient(at 40% 20%, hsla(160, 84%, 85%, 0.3) 0px, transparent 50%), radial-gradient(at 80% 0%, hsla(142, 71%, 90%, 0.2) 0px, transparent 50%), radial-gradient(at 0% 50%, hsla(160, 84%, 95%, 0.3) 0px, transparent 50%)',
            },
            borderRadius: {
                'card': '24px',
            },
            boxShadow: {
                'card': '0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03)',
                'card-hover': '0 20px 25px -5px rgba(0, 0, 0, 0.08), 0 10px 10px -5px rgba(0, 0, 0, 0.02)',
                'soft': '0 2px 15px -3px rgba(0, 0, 0, 0.07), 0 10px 20px -2px rgba(0, 0, 0, 0.04)',
            },
        },
    },
    plugins: [],
}