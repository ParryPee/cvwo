import React from "react";
import TextField from "@mui/material/TextField";
import IconButton from "@mui/material/IconButton";
import Box from "@mui/material/Box";
import SendIcon from "@mui/icons-material/Send";

function TextBox(props: {
	onSubmit: (text: string) => void;
	label: string;
	content?: string;
}) {
	const [text, setText] = React.useState(props.content || "");
	const handleChange = (event: React.ChangeEvent<HTMLTextAreaElement>) => {
		setText(event.target.value);
	};

	const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
		event.preventDefault();
		props.onSubmit(text);
		setText("");
	};

	const handleKeyDown = (event: React.KeyboardEvent<HTMLDivElement>) => {
		if (event.key === "Enter" && !event.shiftKey) {
			event.preventDefault();
			props.onSubmit(text);
			setText("");
		}
	};
	return (
		<Box
			component="form"
			onSubmit={handleSubmit}
			sx={{
				display: "flex",
				flexDirection: "row",
				gap: 2,
				maxWidth: 800,
				marginBottom: 2,
			}}
		>
			<TextField
				label={props.label || "Enter text"}
				value={text}
				onChange={handleChange}
				onKeyDown={handleKeyDown}
				multiline
				maxRows={4}
				size="small"
				sx={{ borderRadius: 2, flexGrow: 1 }}
			/>
			<IconButton
				type="submit"
				color="primary"
				sx={{
					alignSelf: "flex-start",
					backgroundColor: "var(--color-olive-leaf-300)",
				}}
			>
				<SendIcon />
			</IconButton>
		</Box>
	);
}
export default TextBox;
