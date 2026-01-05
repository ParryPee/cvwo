import { useEffect, useState } from "react";
import {
	useNavigate,
	Link as RouterLink,
	useSearchParams,
} from "react-router-dom";
import { fetchAllTopics, searchPost } from "../api/forum"; // Import both
import type { Topic, Post } from "../types/models";
import { timeAgo } from "../utils/date";
import Box from "@mui/material/Box";
import {
	Container,
	CircularProgress,
	Typography,
	Grid,
	CardActionArea,
	Fade,
	Button,
} from "@mui/material";
import { useAuth } from "../context/AuthContext";

const HomePage = () => {
	const navigate = useNavigate();

	const [searchParams] = useSearchParams();
	const query = searchParams.get("q");

	const [topics, setTopics] = useState<Topic[]>([]);
	const [searchResults, setSearchResults] = useState<Post[]>([]);
	const [loading, setLoading] = useState(true);
	const { isAuthenticated, user } = useAuth();

	useEffect(() => {
		const loadData = async () => {
			setLoading(true);
			try {
				if (query) {
					const data = await searchPost(query);
					setSearchResults(data);
				} else {
					const data = await fetchAllTopics();
					setTopics(data);
				}
			} catch (error) {
				console.error("Failed to fetch data:", error);
			} finally {
				setLoading(false);
			}
		};

		loadData();
	}, [query]);

	if (loading) {
		return (
			<Container>
				<CircularProgress />
			</Container>
		);
	} else {
		return (
			<Box
				sx={{
					border: "1px solid #ccc",
					borderRadius: "8px",
					padding: "16px",
					paddingX: "4rem",
					marginTop: "16px",
				}}
			>
				<Box
					sx={{
						display: "flex",
						alignItems: "center",
						justifyContent: "space-between",
						mb: 2,
					}}
				>
					<Typography variant="h4" gutterBottom sx={{ mb: 0 }}>
						{query ? `Results for : ${query}` : "Home Page"}
					</Typography>
					{query ? (
						<Button
							variant="contained"
							onClick={() => navigate("/")}
						>
							Back to Home
						</Button>
					) : (
						isAuthenticated && (
							<Button
								variant="contained"
								color="primary"
								onClick={() => navigate(`/topics/create`)}
							>
								Create Topic
							</Button>
						)
					)}
				</Box>
				{query ? (
					<Grid container spacing={2} sx={{ padding: "16px" }}>
						{!searchResults || searchResults.length === 0 ? (
							<Box
								sx={{
									textAlign: "center",
									padding: "48px 16px",
								}}
							>
								<Typography variant="h6" color="text.secondary">
									No results for "{query}"
								</Typography>
							</Box>
						) : (
							searchResults.map((result, index) => (
								<Fade
									key={result.id}
									in={true}
									timeout={800}
									style={{
										transitionDelay: `${index * 150}ms`,
									}}
								>
									<Grid size={{ xs: 12, sm: 6 }}>
										<CardActionArea
											sx={{
												padding: "16px",
												border: "1px solid #ddd",
												borderRadius: "8px",
												display: "block",
											}}
											component={RouterLink}
											to={`/topics/${result.topic_id}/posts/${result.id}`}
										>
											<Typography
												variant="h6"
												component="h3"
												gutterBottom
												sx={{
													fontWeight: 600,
												}}
											>
												{result.title}
											</Typography>
											<Typography
												variant="body2"
												color="text.secondary"
												gutterBottom
												fontWeight={600}
											>
												Created{" "}
												{timeAgo(result.created_at)}
											</Typography>
											<Typography
												variant="body2"
												color="text.secondary"
												sx={{
													display: "-webkit-box",
													WebkitLineClamp: 2,
													WebkitBoxOrient: "vertical",
													overflow: "hidden",
													textOverflow: "ellipsis",
												}}
											>
												{result.content}
											</Typography>
										</CardActionArea>
									</Grid>
								</Fade>
							))
						)}
					</Grid>
				) : !topics || topics.length === 0 ? (
					<Box
						sx={{
							textAlign: "center",
							padding: "48px 16px",
						}}
					>
						<Typography variant="h6" color="text.secondary">
							No topics available yet
						</Typography>
						<Typography
							variant="body2"
							color="text.secondary"
							sx={{ mt: 1 }}
						>
							Be the first to create a topic!
						</Typography>
					</Box>
				) : (
					<Grid container spacing={2} sx={{ padding: "16px" }}>
						{topics.map((topic, index) => (
							<Fade
								key={topic.id}
								in={true}
								timeout={800}
								style={{ transitionDelay: `${index * 150}ms` }}
							>
								<Grid size={{ xs: 12, sm: 6 }}>
									<CardActionArea
										sx={{
											padding: "16px",
											border: "1px solid #ddd",
											borderRadius: "8px",
											display: "block",
										}}
										component={RouterLink}
										to={`/topics/${topic.id}`}
									>
										<Typography
											variant="h6"
											component="h3"
											gutterBottom
											sx={{
												fontWeight: 600,
											}}
										>
											{topic.title}
										</Typography>
										<Typography
											variant="body2"
											color="text.secondary"
											gutterBottom
											fontWeight={600}
										>
											Created {timeAgo(topic.created_at)}
										</Typography>
										<Typography
											variant="body2"
											color="text.secondary"
											sx={{
												display: "-webkit-box",
												WebkitLineClamp: 2,
												WebkitBoxOrient: "vertical",
												overflow: "hidden",
												textOverflow: "ellipsis",
											}}
										>
											{topic.description}
										</Typography>
									</CardActionArea>
								</Grid>
							</Fade>
						))}
					</Grid>
				)}
			</Box>
		);
	}
};

export default HomePage;
