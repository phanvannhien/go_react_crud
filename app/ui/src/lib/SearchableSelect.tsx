import { useState, useEffect, useRef } from 'react';

export interface SearchableOption {
    value: string;
    label: string;
    sublabel?: string;
}

interface SearchableSelectProps {
    id?: string;
    options: SearchableOption[];
    value: string;
    onChange: (value: string) => void;
    placeholder?: string;
    loading?: boolean;
    onSearch?: (query: string) => void;
    error?: string;
}

export default function SearchableSelect({
    id,
    options,
    value,
    onChange,
    placeholder = 'Search...',
    loading = false,
    onSearch,
    error,
}: SearchableSelectProps) {
    const [isOpen, setIsOpen] = useState(false);
    const [search, setSearch] = useState('');
    const containerRef = useRef<HTMLDivElement>(null);
    const inputRef = useRef<HTMLInputElement>(null);

    // Find selected option label
    const selectedOption = options.find(o => o.value === value);

    // Filter options by search query
    const filtered = search
        ? options.filter(o =>
            o.label.toLowerCase().includes(search.toLowerCase()) ||
            o.value.toLowerCase().includes(search.toLowerCase())
        )
        : options;

    // Close on outside click
    useEffect(() => {
        const handleClick = (e: MouseEvent) => {
            if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
                setIsOpen(false);
                setSearch('');
            }
        };
        document.addEventListener('mousedown', handleClick);
        return () => document.removeEventListener('mousedown', handleClick);
    }, []);

    // Trigger external search
    useEffect(() => {
        if (onSearch) {
            const timer = setTimeout(() => onSearch(search), 300);
            return () => clearTimeout(timer);
        }
    }, [search, onSearch]);

    const handleSelect = (optionValue: string) => {
        onChange(optionValue);
        setIsOpen(false);
        setSearch('');
    };

    const handleOpen = () => {
        setIsOpen(true);
        setTimeout(() => inputRef.current?.focus(), 0);
    };

    return (
        <div ref={containerRef} className="relative">
            {/* Trigger */}
            <button
                type="button"
                id={id}
                onClick={handleOpen}
                className={`w-full px-3 py-2 bg-[var(--color-bg)] border rounded-lg text-sm text-left focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)] transition-colors ${error ? 'border-red-500/50' : 'border-[var(--color-border)]'
                    }`}
            >
                {selectedOption ? (
                    <span className="text-[var(--color-text)]">{selectedOption.label}</span>
                ) : (
                    <span className="text-[var(--color-text-muted)]">{placeholder}</span>
                )}
            </button>

            {/* Dropdown */}
            {isOpen && (
                <div className="absolute z-50 mt-1 w-full bg-[var(--color-surface)] border border-[var(--color-border)] rounded-lg shadow-xl overflow-hidden">
                    {/* Search input */}
                    <div className="p-2 border-b border-[var(--color-border)]">
                        <input
                            ref={inputRef}
                            type="text"
                            value={search}
                            onChange={e => setSearch(e.target.value)}
                            placeholder="Type to search..."
                            className="w-full px-3 py-1.5 bg-[var(--color-bg)] border border-[var(--color-border)] rounded-md text-sm text-[var(--color-text)] focus:outline-none focus:ring-1 focus:ring-[var(--color-primary)]"
                        />
                    </div>

                    {/* Options list */}
                    <div className="max-h-48 overflow-y-auto">
                        {loading && (
                            <div className="px-3 py-3 text-xs text-[var(--color-text-muted)] text-center">Loading...</div>
                        )}
                        {!loading && filtered.length === 0 && (
                            <div className="px-3 py-3 text-xs text-[var(--color-text-muted)] text-center">No results</div>
                        )}
                        {!loading && filtered.map(option => (
                            <button
                                key={option.value}
                                type="button"
                                onClick={() => handleSelect(option.value)}
                                className={`w-full px-3 py-2 text-left text-sm hover:bg-[var(--color-surface-hover)] transition-colors flex items-center justify-between ${option.value === value ? 'bg-[var(--color-primary)]/10 text-[var(--color-primary)]' : 'text-[var(--color-text)]'
                                    }`}
                            >
                                <div>
                                    <span className="block">{option.label}</span>
                                    {option.sublabel && (
                                        <span className="block text-xs text-[var(--color-text-muted)]">{option.sublabel}</span>
                                    )}
                                </div>
                                {option.value === value && (
                                    <span className="text-[var(--color-primary)]">✓</span>
                                )}
                            </button>
                        ))}
                    </div>
                </div>
            )}

            {error && <p className="mt-1 text-xs text-red-400">{error}</p>}
        </div>
    );
}
