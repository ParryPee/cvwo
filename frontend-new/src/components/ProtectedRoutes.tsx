import { Navigate, Outlet, useLocation } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { CircularProgress, Box } from "@mui/material";

const ProtectedRoutes = () => {
	const { isAuthenticated, isLoading } = useAuth();
	const location = useLocation();
	if (isLoading) {
		return (
			<Box display="flex" justifyContent="center" mt={4}>
				<CircularProgress />
			</Box>
		);
	}
	return isAuthenticated ? (
		<Outlet />
	) : (
		<Navigate to="/login" replace state={{ from: location }} />
	);
};

export default ProtectedRoutes;
