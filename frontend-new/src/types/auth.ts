import type { User } from "./models";

interface AuthContextType {
    user: User | null;
    isAuthenticated: boolean;
    isLoading: boolean;
    login: (username: string) => Promise<void>;
    logout: () => void;
}

export type { AuthContextType };