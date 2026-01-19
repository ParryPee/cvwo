import { useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { fetchAllPosts, fetchAllTopics, searchGlobal } from "../api/forum";
import type { SearchResult, Topic } from "../types/models";
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
	Card,
	CardContent,
	CardActionArea,
	Stack,
} from "@mui/material";
import AddIcon from "@mui/icons-material/Add";
import ClearIcon from "@mui/icons-material/Clear";
import WhatshotIcon from "@mui/icons-material/Whatshot";
import ForumIcon from "@mui/icons-material/Forum";
import { useAuth } from "../context/AuthContext";
import Feed from "../components/Feed";

const HomePage = () => {
	const navigate = useNavigate();
	const [searchParams] = useSearchParams();
	const query = searchParams.get("q");

	const [data, setData] = useState<SearchResult>({ posts: [], topics: [] });
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
					const data = await searchGlobal(query);
					setData(data);
					setSearchTerm(query);
					if (!data.topics) {
						data.topics = [];
					}
					if (!data.posts) {
						data.posts = [];
					}
				} else {
					const postData = await fetchAllPosts(LIMIT, 0);
					const topicData: Topic[] = await fetchAllTopics();
					setData({ posts: postData, topics: topicData });
				}
			} catch (error) {
				console.error("Failed to load data", error);
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
	console.log("HomePage data:", data);
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
						Search for posts! Or select a topic to get started.
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

						{loading ? (
							<Box
								sx={{
									display: "flex",
									justifyContent: "center",
									mt: 4,
								}}
							>
								<CircularProgress />
							</Box>
						) : (
							<>
								{query && data.topics.length > 0 && (
									<Box sx={{ mb: 4 }}>
										<Typography
											variant="h6"
											gutterBottom
											sx={{
												display: "flex",
												alignItems: "center",
												gap: 1,
											}}
										>
											<ForumIcon
												color="primary"
												fontSize="small"
											/>
											Communities
										</Typography>
										<Stack
											direction="row"
											spacing={2}
											sx={{
												overflowX: "auto",
												pb: 1,
											}}
										>
											{data.topics.map((topic) => (
												<Card
													key={topic.id}
													variant="outlined"
													sx={{
														minWidth: 200,
														maxWidth: 250,
														flexShrink: 0,
													}}
												>
													<CardActionArea
														onClick={() =>
															navigate(
																`/topics/${topic.id}`,
															)
														}
													>
														<CardContent
															sx={{ py: 2 }}
														>
															<Typography
																variant="subtitle1"
																fontWeight="bold"
																noWrap
															>
																{topic.title}
															</Typography>
															<Typography
																variant="caption"
																color="text.secondary"
																noWrap
																display="block"
															>
																{
																	topic.description
																}
															</Typography>
														</CardContent>
													</CardActionArea>
												</Card>
											))}
										</Stack>
										<Divider sx={{ mt: 3 }} />
									</Box>
								)}
							</>
						)}
						{query &&
							data.topics.length > 0 &&
							data.posts.length > 0 && (
								<Typography
									variant="h6"
									gutterBottom
									sx={{ mt: 2 }}
								>
									Posts
								</Typography>
							)}
						<Feed
							items={data.posts}
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
						<Box>
							<Paper
								sx={{ p: 3, borderRadius: 2 }}
								variant="outlined"
							>
								<Typography
									variant="h6"
									gutterBottom
									sx={{ mb: 2 }}
								>
									Top Communities
								</Typography>
								{data.topics
									.sort((a, b) => b.post_count - a.post_count)
									.slice(0, 3)
									.map((topic) => (
										<CardActionArea
											key={topic.id}
											sx={{
												mb: 2,
												cursor: "pointer",
											}}
											onClick={() =>
												navigate(`/topics/${topic.id}`)
											}
										>
											<Typography
												variant="subtitle1"
												fontWeight="bold"
												noWrap
											>
												{topic.title}
											</Typography>
											<Typography
												variant="caption"
												color="text.secondary"
												noWrap
												display="block"
											>
												{topic.post_count} posts
											</Typography>
										</CardActionArea>
									))}
							</Paper>
						</Box>
					</Grid>
				</Grid>
			</Container>
		</Box>
	);
};

export default HomePage;
