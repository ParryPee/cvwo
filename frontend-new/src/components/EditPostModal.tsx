import { useState, useEffect } from "react";
import {
	Dialog,
	DialogTitle,
	DialogContent,
	DialogActions,
	TextField,
	Button,
} from "@mui/material";

interface EditPostModalProps {
	open: boolean;
	onClose: () => void;
	onSubmit: (newTitle: string, newContent: string) => Promise<void>;
	initialTitle: string;
	initialContent: string;
}

const EditPostModal = ({
	open,
	onClose,
	onSubmit,
	initialTitle,
	initialContent,
}: EditPostModalProps) => {
	const [title, setTitle] = useState(initialTitle);
	const [content, setContent] = useState(initialContent);
	const [loading, setLoading] = useState(false);

	useEffect(() => {
		if (open) {
			setContent(initialContent);
		}
	}, [open, initialContent]);

	const handleSubmit = async () => {
		if (content.trim() === "") {
			return;
		}
		try {
			setLoading(true);
			await onSubmit(title, content);
		} catch (error) {
			console.error("Failed to update post", error);
		} finally {
			setLoading(false);
		}
	};
	return (
		<Dialog open={open} onClose={onClose} fullWidth maxWidth="sm">
			<DialogTitle>Edit Post</DialogTitle>
			<DialogContent>
				<TextField
					autoFocus
					margin="dense"
					label="Content"
					fullWidth
					multiline
					rows={6}
					variant="outlined"
					value={title}
					onChange={(e) => setTitle(e.target.value)}
					disabled={loading}
				/>
				<TextField
					autoFocus
					margin="dense"
					label="Content"
					fullWidth
					multiline
					rows={6}
					variant="outlined"
					value={content}
					onChange={(e) => setContent(e.target.value)}
					disabled={loading}
				/>
			</DialogContent>
			<DialogActions>
				<Button onClick={onClose} color="secondary" disabled={loading}>
					Cancel
				</Button>
				<Button
					onClick={handleSubmit}
					variant="contained"
					color="primary"
					disabled={loading || !content.trim()}
				>
					{loading ? "Saving..." : "Save"}
				</Button>
			</DialogActions>
		</Dialog>
	);
};

export default EditPostModal;
