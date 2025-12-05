export default function Card({ children, className = '', padding = 'p-6' }) {
    return (
        <div className={`bg-white rounded-card shadow-lg ${padding} ${className}`}>
            {children}
        </div>
    );
}