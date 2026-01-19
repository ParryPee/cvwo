import React, { useState } from "react";
import { useNavigate, Link as RouterLink } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import {
	Container,
	Box,
	Typography,
	TextField,
	Button,
	Alert,
	Paper,
	Link,
} from "@mui/material";
import { AxiosError } from "axios";

const LoginPage: React.FC = () => {
	const { login } = useAuth();

	const navigate = useNavigate();

	const [username, setUsername] = useState("");
	const [error, setError] = useState("");
	const [loading, setLoading] = useState(false);

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault(); // Stop the page from reloading
		setError("");
		setLoading(true);

		try {
			await login(username);

			navigate("/");
		} catch (err) {
			console.error(err);
			setError(
				err instanceof AxiosError ? err.response?.data : "Login failed",
			);
		} finally {
			setLoading(false);
		}
	};

	return (
		<Container
			maxWidth={false}
			disableGutters
			sx={{
				backgroundImage:
					"linear-gradient(to bottom, var(--color-lavender-grey-700), var(--color-space-indigo-500))",
			}}
		>
			<Box
				sx={{
					display: "flex",
					flexDirection: "column",
					alignItems: "center",
					padding: "32px",
				}}
				maxWidth="700px"
				margin="0 auto"
			>
				<Paper
					elevation={3}
					sx={{
						padding: 4,
						width: "100%",
						bgcolor: "var(--color-platinum-100)",
						borderRadius: 4,
					}}
				>
					<Typography
						component="h1"
						variant="h5"
						align="center"
						gutterBottom
					>
						Welcome Back
					</Typography>

					{error && (
						<Alert severity="error" sx={{ mb: 2 }}>
							{error}
						</Alert>
					)}

					<Box component="form" onSubmit={handleSubmit}>
						<TextField
							margin="normal"
							required
							fullWidth
							id="username"
							label="Username"
							name="username"
							autoFocus
							value={username}
							onChange={(e) => setUsername(e.target.value)}
							disabled={loading}
						/>

						<Button
							type="submit"
							variant="contained"
							sx={{
								mt: 3,
								mb: 2,
								borderRadius: 4,
								bgcolor: "var(--color-flag-red-500)",
								":hover": {
									backgroundColor:
										"var(--color-flag-red-600)",
								},
								color: "white",
								margin: "0 auto",
								display: "block",
								maxWidth: "150px",
							}}
							disabled={loading}
						>
							{loading ? "Signing in..." : "Sign In"}
						</Button>
					</Box>

					<Typography variant="body2" align="center" sx={{ mt: 2 }}>
						Don&apos;t have an account?{" "}
						<Link component={RouterLink} to="/register">
							Register
						</Link>
					</Typography>
				</Paper>
			</Box>
		</Container>
	);
};

export default LoginPage;
