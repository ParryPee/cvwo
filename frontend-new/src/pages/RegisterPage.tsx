import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { register } from "../api/auth";
import {
	Container,
	Box,
	Typography,
	TextField,
	Button,
	Alert,
	Paper,
} from "@mui/material";
import { AxiosError } from "axios";

const LoginPage: React.FC = () => {
	const { isAuthenticated, login } = useAuth();

	const navigate = useNavigate();

	const [username, setUsername] = useState("");
	const [error, setError] = useState("");
	const [loading, setLoading] = useState(false);

	if (isAuthenticated) {
		navigate("/");
	}

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault(); // Stop the page from reloading
		setError("");
		setLoading(true);

		try {
			await register(username);
			await login(username);

			navigate("/");
		} catch (err) {
			console.error(err);
			setError(
				err instanceof AxiosError
					? err.response?.data
					: "Registration failed",
			);
		} finally {
			setLoading(false);
		}
	};

	return (
		<Container
			maxWidth={false}
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
						Join the Community!
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
							fullWidth
							variant="contained"
							sx={{
								mb: 4,
								backgroundColor: "var(--color-flag-red-500)",
								":hover": {
									backgroundColor:
										"var(--color-flag-red-600)",
								},
								borderRadius: 4,
								maxWidth: "150px",
								margin: "20px auto",
								display: "block",
							}}
							disabled={loading}
						>
							{loading ? "Registering" : "Register"}
						</Button>
					</Box>

					<Button
						type="button"
						variant="outlined"
						fullWidth
						sx={{
							backgroundColor: "var(--color-flag-red-500)",
							":hover": {
								backgroundColor: "var(--color-flag-red-600)",
							},
							borderRadius: 4,
							maxWidth: "150px",
							margin: "0 auto",
							display: "block",
							color: "white",
						}}
						onClick={() => navigate("/login")}
					>
						Login
					</Button>
				</Paper>
			</Box>
		</Container>
	);
};

export default LoginPage;
