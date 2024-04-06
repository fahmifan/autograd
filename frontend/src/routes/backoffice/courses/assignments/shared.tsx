import { Paper } from "@mantine/core";
import {
	MDXEditor,
	MDXEditorMethods,
	headingsPlugin,
	listsPlugin,
	markdownShortcutPlugin,
	quotePlugin,
	thematicBreakPlugin,
} from "@mdxeditor/editor";
import { forwardRef } from "react";

type MarkdownEditorProps = {
	defaultValue?: string;
	onChange?: (value: string) => void;
	ref: React.RefObject<MDXEditorMethods>;
};

export const MarkdownEditor = forwardRef<MDXEditorMethods, MarkdownEditorProps>(
	(props: MarkdownEditorProps, ref) => {
		return (
			<Paper shadow="xs" p="xl">
				<MDXEditor
					ref={ref}
					onChange={props.onChange}
					markdown={props.defaultValue || "## Description"}
					plugins={[
						headingsPlugin(),
						listsPlugin(),
						quotePlugin(),
						thematicBreakPlugin(),
						markdownShortcutPlugin(),
					]}
				/>
			</Paper>
		);
	},
);
