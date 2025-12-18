import React, { createContext, useContext, useState, useEffect } from "react";
import type { AuthContextType } from "../types/auth";
import {
	login as apiLogin,
	fetchCurrentUser,
	logout as apiLogout,
} from "../api/auth";
import type { User } from "../types/models";

const defaultContextValue: AuthContextType = {
	user: null,
	isAuthenticated: false,
	isLoading: true,
	login: async (username: string) => {},
	logout: () => {},
};

const AuthContext = createContext<AuthContextType>(defaultContextValue);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
	children,
}) => {
	const [user, setUser] = useState<User | null>(null);
	const [isLoading, setIsLoading] = useState<boolean>(true);
	const isAuthenticated = !!user;
	useEffect(() => {
		const initializeAuth = async () => {
			try {
				const userData = await fetchCurrentUser();
				setUser(userData);
			} catch (error) {
				console.log("No valid session found");
				setUser(null);
			} finally {
				setIsLoading(false);
			}
		};

		initializeAuth();
	}, []);
	const login = async (username: string) => {
		await apiLogin(username);
		const userData = await fetchCurrentUser();
		setUser(userData);
	};
	const logout = async () => {
		try {
			await apiLogout();
		} catch (error) {
			console.error("Logout failed:", error);
		} finally {
			setUser(null);
		}
	};

	return (
		<AuthContext.Provider
			value={{ user, isAuthenticated, isLoading, login, logout }}
		>
			{children}
		</AuthContext.Provider>
	);
};
export const useAuth = () => {
	const context = useContext(AuthContext);
	if (!context)
		throw new Error("useAuth must be used within an AuthProvider");
	return context;
};
