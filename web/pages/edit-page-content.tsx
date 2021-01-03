import React, { useEffect } from 'react';

import MarkdownIt from 'markdown-it';
const mdParser = new MarkdownIt(/* Markdown-it options */);

import dynamic from 'next/dynamic';
import 'react-markdown-editor-lite/lib/index.css';

import { SERVER_CONFIG, fetchData, FETCH_INTERVAL, UPDATE_CHAT_MESSGAE_VIZ } from "../utils/apis";

import { Table, Typography, Tooltip, Button } from "antd";

const MdEditor = dynamic(() => import('react-markdown-editor-lite'), {
    ssr: false
});

export default function PageContentEditor() {
    const [content, setContent] = React.useState("");

    function handleEditorChange({ html, text }) {
        setContent(text);
    }

    function handleSave() {
        console.log(content);
        alert("Make API call to save here." + content)
    }

    async function setInitialContent() {
        const serverConfig = await fetchData(SERVER_CONFIG);
        const initialContent = serverConfig.instanceDetails.extraPageContent;
        setContent(initialContent);
    }

    useEffect(() => {
        setInitialContent();
    }, []);

    return (
        <div>
            <MdEditor
                style={{ height: "500px" }}
                value={content}
                renderHTML={(content) => mdParser.render(content)}
                onChange={handleEditorChange}
                config={{ htmlClass: 'markdown-editor-preview-pane', markdownClass: 'markdown-editor-pane' }}
            />
            <Button onClick={handleSave}>Save</Button>
        </div>
    );
}
