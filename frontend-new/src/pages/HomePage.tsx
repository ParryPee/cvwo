import { useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { fetchAllPosts, searchPost } from "../api/forum"; // <--- Updated import
import type { Post } from "../types/models";
import {
	Container,
	Grid,
	Box,
	Typography,
	Button,
	TextField,
	Paper,
	Divider,
	CircularProgress,
} from "@mui/material";
import AddIcon from "@mui/icons-material/Add";
import ClearIcon from "@mui/icons-material/Clear";
import WhatshotIcon from "@mui/icons-material/Whatshot";
import { useAuth } from "../context/AuthContext";
import Feed from "../components/Feed";

const HomePage = () => {
	const navigate = useNavigate();
	const [searchParams] = useSearchParams();
	const query = searchParams.get("q");

	const [posts, setPosts] = useState<Post[]>([]);
	const [loading, setLoading] = useState(true);
	const [searchTerm, setSearchTerm] = useState("");

	const LIMIT = 10;
	const { isAuthenticated } = useAuth();

	const handleSearch = (e: React.FormEvent) => {
		e.preventDefault();
		if (searchTerm.trim()) navigate(`/?q=${searchTerm}`);
	};

	const handleClearSearch = () => {
		setSearchTerm("");
		navigate("/");
	};

	useEffect(() => {
		const loadData = async () => {
			setLoading(true);
			try {
				if (query) {
					const data = await searchPost(query);
					setPosts(data);
					setSearchTerm(query);
				} else {
					const data = await fetchAllPosts(LIMIT, 0);
					setPosts(data);
				}
			} catch (error) {
				console.error("Failed to fetch posts:", error);
			} finally {
				setLoading(false);
			}
		};

		loadData();
	}, [query]);

	if (loading) {
		return (
			<Container
				sx={{ mt: 8, display: "flex", justifyContent: "center" }}
			>
				<CircularProgress />
			</Container>
		);
	}

	return (
		<Box sx={{ pb: 8 }}>
			<Paper
				elevation={0}
				sx={{
					bgcolor: "primary.main",
					color: "primary.contrastText",
					pt: 8,
					pb: 8,
					mb: 4,
					borderRadius: 0,
				}}
			>
				<Container maxWidth="md" sx={{ textAlign: "center" }}>
					<Typography
						variant="h3"
						component="h1"
						fontWeight="bold"
						gutterBottom
					>
						Welcome to the Community
					</Typography>
					<Typography variant="h6" sx={{ mb: 4, opacity: 0.9 }}>
						Search for posts!
					</Typography>

					<Paper
						component="form"
						onSubmit={handleSearch}
						sx={{
							p: "2px 4px",
							display: "flex",
							alignItems: "center",
							width: "100%",
							maxWidth: 600,
							mx: "auto",
							borderRadius: "50px",
							pl: 2,
						}}
					>
						<TextField
							fullWidth
							placeholder="Search for posts..."
							variant="standard"
							value={searchTerm}
							onChange={(e) => setSearchTerm(e.target.value)}
						/>
						<Button
							type="submit"
							variant="contained"
							sx={{ borderRadius: "50px", px: 4, m: 0.5 }}
						>
							Search
						</Button>
					</Paper>
				</Container>
			</Paper>

			<Container maxWidth="lg">
				<Grid container spacing={4}>
					<Grid size={{ xs: 12, md: 8 }}>
						<Box
							sx={{
								display: "flex",
								justifyContent: "space-between",
								alignItems: "center",
								mb: 2,
							}}
						>
							{query ? (
								<Box
									sx={{
										display: "flex",
										alignItems: "center",
										gap: 2,
									}}
								>
									<Typography variant="h5" fontWeight="bold">
										Results for "{query}"
									</Typography>
									<Button
										startIcon={<ClearIcon />}
										onClick={handleClearSearch}
										size="small"
										color="inherit"
									>
										Clear
									</Button>
								</Box>
							) : (
								<Box
									sx={{
										display: "flex",
										alignItems: "center",
										gap: 1,
									}}
								>
									<WhatshotIcon color="error" />
									<Typography variant="h5" fontWeight="bold">
										Latest Activity
									</Typography>
								</Box>
							)}

							{isAuthenticated && (
								<Button
									variant="outlined"
									startIcon={<AddIcon />}
									onClick={() => navigate(`/topics/create`)}
									size="small"
								>
									New Topic
								</Button>
							)}
						</Box>

						<Divider sx={{ mb: 3 }} />

						<Feed
							items={posts}
							emptyMessage={
								query
									? `No posts found for "${query}"`
									: "No posts yet. Be the first to start a conversation!"
							}
						/>
					</Grid>

					<Grid size={{ xs: 12, md: 4 }}>
						<Box sx={{ position: "sticky", top: 24 }}>
							<Paper
								sx={{ p: 3, mb: 3, borderRadius: 2 }}
								variant="outlined"
							>
								<Typography
									variant="body2"
									color="text.secondary"
								>
									Welcome to the forum, join in the discussion
									or start your own!
								</Typography>
								<Button
									fullWidth
									variant="contained"
									onClick={() => navigate("/topics/create")}
								>
									Start your own topic!
								</Button>
							</Paper>
						</Box>
					</Grid>
				</Grid>
			</Container>
		</Box>
	);
};

export default HomePage;
