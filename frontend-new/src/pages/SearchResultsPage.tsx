import { useEffect, useState } from "react";
import {
	useNavigate,
	Link as RouterLink,
	useSearchParams,
} from "react-router-dom";
import { searchPost } from "../api/forum";
import type { Post } from "../types/models";
import { timeAgo, formatDate } from "../utils/date";
import Box from "@mui/material/Box";
import {
	Container,
	CircularProgress,
	Typography,
	Grid,
	Card,
	Button,
	CardActionArea,
	Fade,
} from "@mui/material";
import { Link, Padding } from "@mui/icons-material";
import { useAuth } from "../context/AuthContext";

const HomePage = () => {
	const navigate = useNavigate();
	const [posts, setPosts] = useState<Post[]>([]);
	const [loading, setLoading] = useState(true);
	const [searchParams] = useSearchParams();
	const { isAuthenticated } = useAuth();

	const q = searchParams.get("q");
	useEffect(() => {
		const loadResults = async () => {
			if (!q) {
				console.error("Invalid search query");
				return;
			}
			try {
				const data = await searchPost(q);
				setPosts(data);
			} catch (error) {
				console.error("Failed to fetch topics:", error);
			} finally {
				setLoading(false);
			}
		};

		loadResults();
	}, []);

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
						Search Results
					</Typography>
					{isAuthenticated && (
						<Button
							variant="contained"
							color="primary"
							onClick={() => navigate(`/topics/create`)}
						>
							Can't find what you're looking for? Start your own
							topic!
						</Button>
					)}
				</Box>
				{!posts || posts.length === 0 ? (
					<Box
						sx={{
							textAlign: "center",
							padding: "48px 16px",
						}}
					>
						<Typography variant="h6" color="text.secondary">
							No results found. Try another query.
						</Typography>
					</Box>
				) : (
					<Grid container spacing={2} sx={{ padding: "16px" }}>
						{posts.map((post, index) => (
							<Fade
								key={post.id}
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
										to={`/topics/${post.topic_id}/posts/${post.id}`}
									>
										<Typography
											variant="h6"
											component="h3"
											gutterBottom
											sx={{
												fontWeight: 600,
											}}
										>
											{post.title}
										</Typography>
										<Typography
											variant="body2"
											color="text.secondary"
											gutterBottom
											fontWeight={600}
										>
											Created {timeAgo(post.created_at)}
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
											{post.content}
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
