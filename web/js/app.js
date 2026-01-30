let workspace = null;
let currentTatuCode = '';

// init the blocks and workspace
function initWorkspace() {
    const config = {
        toolbox: {
            kind: 'categoryToolbox',

            contents: [
                {
                    "kind": "category",
                    "name": "Variables",
                    "colour": "0",
                    "contents": [
                        {
                            "kind": "block",
                            "type": "variables_create_tatu"
                        },
                        {
                            "kind": "block",
                            "type": "variables_set"
                        },
                        {
                            "kind": "block",
                            "type": "variables_get"
                        }
                    ]
                },
                {
                    "kind": "category",
                    "name": "Control",
                    "colour": "30",
                    "contents": [
                        {
                            "kind": "block",
                            "type": "controls_if"
                        },
                        {
                            "kind": "block",
                            "type": "controls_repeat_ext"
                        },
                        {
                            "kind": "block",
                            "type": "controls_whileUntil"
                        },
                        {
                            "kind": "block",
                            "type": "controls_for"
                        },
                        {
                            "kind": "block",
                            "type": "controls_forEach"
                        },
                        {
                            "kind": "block",
                            "type": "controls_flow_statements"
                        },
                        {
                            "kind": "block",
                            "type": "expression_tatu"
                        }
                    ]
                },
                {
                    "kind": "category",
                    "name": "Logic",
                    "colour": "60",
                    "contents": [
                        {
                            "kind": "block",
                            "type": "logic_boolean"
                        },
                        {
                            "kind": "block",
                            "type": "logic_null"
                        },
                        {
                            "kind": "block",
                            "type": "logic_compare"
                        },
                        {
                            "kind": "block",
                            "type": "logic_operation"
                        },
                        {
                            "kind": "block",
                            "type": "logic_negate"
                        },
                        {
                            "kind": "block",
                            "type": "logic_ternary"
                        },
                        {
                            "kind": "block",
                            "type": "type_check_call_tatu"
                        }
                    ]
                },
                {
                    "kind": "category",
                    "name": "Math",
                    "colour": "90",
                    "contents": [
                        {
                            "kind": "block",
                            "type": "math_number"
                        },
                        {
                            "kind": "block",
                            "type": "math_arithmetic"
                        },
                        {
                            "kind": "block",
                            "type": "math_single"
                        },
                        {
                            "kind": "block",
                            "type": "math_trig"
                        },
                        {
                            "kind": "block",
                            "type": "math_constant"
                        },
                        {
                            "kind": "block",
                            "type": "math_round"
                        },
                        {
                            "kind": "block",
                            "type": "math_modulo"
                        },
                        {
                            "kind": "block",
                            "type": "math_random_int"
                        },
                        {
                            "kind": "block",
                            "type": "math_random_float"
                        },
                        {
                            "kind": "block",
                            "type": "math_constrain"
                        },
                        {
                            "kind": "block",
                            "type": "math_number_property"
                        },
                        {
                            "kind": "block",
                            "type": "math_change"
                        },
                        {
                            "kind": "block",
                            "type": "math_call_tatu"
                        }
                    ]
                },
                {
                    "kind": "category",
                    "name": "Text",
                    "colour": "120",
                    "contents": [
                        {
                            "kind": "block",
                            "type": "text"
                        },
                        {
                            "kind": "block",
                            "type": "text_join"
                        },
                        {
                            "kind": "block",
                            "type": "text_length"
                        },
                        {
                            "kind": "block",
                            "type": "text_isEmpty"
                        },
                        {
                            "kind": "block",
                            "type": "text_charAt"
                        },
                        {
                            "kind": "block",
                            "type": "text_getSubstring"
                        },
                        {
                            "kind": "block",
                            "type": "text_append"
                        },
                        {
                            "kind": "block",
                            "type": "text_indexOf"
                        },
                        {
                            "kind": "block",
                            "type": "text_changeCase"
                        },
                        {
                            "kind": "block",
                            "type": "text_trim"
                        },
                        {
                            "kind": "block",
                            "type": "text_print"
                        },
                        {
                            "kind": "block",
                            "type": "string_call_tatu"
                        }
                    ]
                },
                {
                    "kind": "category",
                    "name": "List (vector)",
                    "colour": "150",
                    "contents": [
                        {
                            "kind": "block",
                            "type": "lists_create_with"
                        },
                        {
                            "kind": "block",
                            "type": "lists_length"
                        },
                        {
                            "kind": "block",
                            "type": "lists_isEmpty"
                        },
                        {
                            "kind": "block",
                            "type": "lists_getIndex"
                        },
                        {
                            "kind": "block",
                            "type": "lists_setIndex"
                        },
                        {
                            "kind": "block",
                            "type": "lists_indexOf"
                        },
                        {
                            "kind": "block",
                            "type": "lists_getSublist"
                        },
                        {
                            "kind": "block",
                            "type": "lists_split"
                        },
                        {
                            "kind": "block",
                            "type": "lists_sort"
                        },
                        {
                            "kind": "block",
                            "type": "lists_repeat"
                        },
                        {
                            "kind": "block",
                            "type": "lists_reverse"
                        },
                        {
                            "kind": "block",
                            "type": "vector_call_tatu"
                        }
                    ]
                },
                {
                    "kind": "category",
                    "name": "Map",
                    "colour": "180",
                    "contents": [
                        {
                            "kind": "block",
                            "type": "map_create_with_tatu"
                        },
                        {
                            "kind": "block",
                            "type": "map_pair_tatu"
                        },
                        {
                            "kind": "block",
                            "type": "map_call_tatu"
                        }
                    ]
                },
                {
                    "kind": "category",
                    "name": "Function",
                    "colour": "210",
                    "contents": [
                        {
                            "kind": "block",
                            "type": "procedures_defnoreturn"
                        },
                        {
                            "kind": "block",
                            "type": "procedures_callnoreturn"
                        },
                        {
                            "kind": "block",
                            "type": "procedures_defreturn"
                        },
                        {
                            "kind": "block",
                            "type": "procedures_ifreturn"
                        },
                        {
                            "kind": "block",
                            "type": "procedures_callreturn"
                        },
                        {
                            "kind": "block",
                            "type": "function_call_tatu"
                        }
                    ]
                },
                {
                    "kind": "category",
                    "name": "Casting",
                    "colour": "240",
                    "contents": [
                        {
                            "kind": "block",
                            "type": "casting_call_tatu"
                        }
                    ]
                },
                {
                    "kind": "category",
                    "name": "Time",
                    "colour": "270",
                    "contents": [
                        {
                            "kind": "block",
                            "type": "time_call_tatu"
                        }
                    ]
                },
                {
                    "kind": "category",
                    "name": "JSON",
                    "colour": "300",
                    "contents": [
                        {
                            "kind": "block",
                            "type": "json_call_tatu"
                        }
                    ]
                },
                {
                    "kind": "category",
                    "name": "Regex",
                    "colour": "330",
                    "contents": [
                        {
                            "kind": "block",
                            "type": "regex_call_tatu"
                        }
                    ]
                },
                {
                    "kind": "category",
                    "name": "File System",
                    "colour": "345",
                    "contents": [
                        {
                            "kind": "block",
                            "type": "filesystem_call_tatu"
                        }
                    ]
                },
                {
                    "kind": "category",
                    "name": "Include",
                    "colour": "360",
                    "contents": [
                        {
                            "kind": "block",
                            "type": "include_file_tatu"
                        }
                    ]
                },
                {
                    "kind": "category",
                    "name": "Comment",
                    "colour": "60",
                    "contents": [
                        {
                            "kind": "block",
                            "type": "comment_tatu"
                        }
                    ]
                }
            ]
        },

        move: {
            scrollbars: {
                horizontal: true,
                vertical: true
            },
            drag: true,
            wheel: false
        },

        zoom: {
            controls: true,
            wheel: true,
            startScale: 1.0,
            maxScale: 3,
            minScale: 0.3,
            scaleSpeed: 1.1,
            pinch: true
        },

        grid: {
            spacing: 20,
            length: 3,
            colour: '#ccc',
            snap: true
        },

        trashcan: true,
    };

    registerBlocks();

    workspace = Blockly.inject('blocklyCanvas', config);
}

// init the import
function initImportJSON() {
    document.getElementById('fileInput').click();
}

// import JSON to workspace
function importJSON(event) {
    const file = event.target.files[0];

    if (file) {
        const reader = new FileReader();

        reader.onload = function(e) {
            try {
                const data = JSON.parse(e.target.result);
                Blockly.serialization.workspaces.load(data.workspace, workspace);
                if (data.inputs) {
                    document.getElementById('inputVars').value = data.inputs;
                }
                clearOutput();
            } catch (error) {
                showError('blocklyCode', 'importing file: ' + error.message);
            }
        };
        reader.readAsText(file);
    }
    event.target.value = '';
}

// import workspace to JSON
function exportJSON() {
    const state = Blockly.serialization.workspaces.save(workspace);
    const inputVars = document.getElementById('inputVars').value;
    const data = {
        workspace: state,
        inputs: inputVars
    };
    const json = JSON.stringify(data, null, 2);
    const blob = new Blob([json], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');

    a.href = url;
    a.download = 'blockly-workspace.json';
    a.click();
    URL.revokeObjectURL(url);
}

// transpile blockly code to Tatu
async function transpile(){
    const state = Blockly.serialization.workspaces.save(workspace);

    try {
        const response = await fetch('/transpile', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(state)
        });

        if (!response.ok) {
            const errorMessage = await response.text();
            currentTatuCode = '';
            document.getElementById('btnExecute').disabled = true;
            showError('blocklyCode', errorMessage);
            return;
        }

        const result = await response.json();
        document.getElementById('blocklyCode').innerHTML =
            '<h3>Blockly</h3>' +
            '<pre class="output-json">' + escapeHtml(result.json) + '</pre>' +
            '<h3>Tatu</h3>' +
            '<pre class="output-tatu"><code class="language-lisp">' + escapeHtml(result.tatu_code) + '</code></pre>';

        hljs.highlightAll();

        currentTatuCode = result.tatu_code;
        document.getElementById('btnExecute').disabled = false;
        document.getElementById('executionOutput').innerHTML = '';
    } catch (error) {
        currentTatuCode = '';
        document.getElementById('btnExecute').disabled = true;
        showError('blocklyCode', error.message);
    }
}

// execute Tatu code
async function execute(){
    if (!currentTatuCode) {
        return;
    }

    const inputText = document.getElementById('inputVars').value.trim();
    let inputs = {};

    if (inputText) {
        try {
            inputs = JSON.parse(inputText);
        } catch (error) {
            showError('executionOutput', 'Invalid JSON in inputs: ' + error.message);
            return;
        }
    }

    try {
        const response = await fetch('/execute', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ code: currentTatuCode, inputs: inputs })
        });

        const output = await response.text();

        if (!response.ok) {
            showError('executionOutput', output);
            return;
        }

        document.getElementById('executionOutput').innerHTML =
            '<pre class="output-execution">' + escapeHtml(output) + '</pre>';
    } catch (error) {
        showError('executionOutput', error.message);
    }
}

// clear output function
function clearOutput() {
    document.getElementById('blocklyCode').innerHTML = '';
    document.getElementById('executionOutput').innerHTML = '';
    document.getElementById('btnExecute').disabled = true;
    currentTatuCode = '';
}

// escape HTML characters
function escapeHtml(str) {
    return str
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#039;');
}

// show error message
function showError(elementId, message) {
    document.getElementById(elementId).innerHTML =
        '<pre class="output-error">Error: ' + escapeHtml(message) + '</pre>';
}

document.addEventListener('DOMContentLoaded', initWorkspace);
document.getElementById('btnImport').addEventListener('click', initImportJSON);
document.getElementById('fileInput').addEventListener('change', importJSON);
document.getElementById('btnExport').addEventListener('click', exportJSON);
document.getElementById('btnTranspile').addEventListener('click', transpile);
document.getElementById('btnExecute').addEventListener('click', execute);
document.getElementById('btnClear').addEventListener('click', clearOutput);
