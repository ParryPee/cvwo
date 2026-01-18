import { Card, CardActionArea, Typography, Box, Link } from "@mui/material";
import { Link as RouterLink } from "react-router-dom";
import { timeAgo } from "../utils/date";

interface DisplayCardProps {
	title: string;
	previewText: string;
	createdAt: string;
	linkTo: string;
	username?: string;
	timeAgo?: string;
	topicTitle?: string;
	topicID?: number;
}

const DisplayCard = ({
	title,
	previewText,
	createdAt,
	linkTo,
	...props
}: DisplayCardProps) => {
	return (
		<Card
			sx={{
				height: "100%",
				display: "flex",
				flexDirection: "column",
				border: "1px solid #ddd",
				borderRadius: "8px",
				transition: "transform 0.2s, box-shadow 0.2s",
				"&:hover": {
					transform: "translateY(-2px)",
					boxShadow: 4,
				},
			}}
		>
			<CardActionArea
				component={RouterLink}
				to={linkTo}
				sx={{
					p: 2,
					height: "100%",
					display: "flex",
					flexDirection: "column",
					alignItems: "flex-start",
					justifyContent: "flex-start",
				}}
			>
				<Typography
					variant="h6"
					component="h3"
					gutterBottom
					sx={{ fontWeight: 600 }}
				>
					{title}
				</Typography>

				<Typography
					variant="body2"
					color="text.secondary"
					gutterBottom
					fontWeight={600}
					sx={{ mb: 1.5 }}
				>
					Created {timeAgo(createdAt)}
				</Typography>

				<Typography
					variant="body2"
					color="text.secondary"
					sx={{
						display: "-webkit-box",
						WebkitLineClamp: 3,
						WebkitBoxOrient: "vertical",
						overflow: "hidden",
						textOverflow: "ellipsis",
					}}
				>
					{previewText}
				</Typography>
			</CardActionArea>
			<Box
				sx={{
					display: "flex",
					flexDirection: "row",
					justifyContent: "center",
				}}
			>
				{props.topicID ? (
					<Link
						variant="caption"
						sx={{ color: "primary.main", mb: 1, display: "block" }}
						href={`/topics/${props.topicID}`}
					>
						Posted in {props.topicTitle || "General"}
					</Link>
				) : (
					<Typography
						variant="caption"
						sx={{ color: "primary.main", mb: 1, display: "block" }}
					>
						Posted in {props.topicTitle || "General"}
					</Typography>
				)}
			</Box>
		</Card>
	);
};
export default DisplayCard;
