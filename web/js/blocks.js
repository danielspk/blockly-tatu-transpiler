// custom tatu call function block
function createFunctionCallBlock(blockType, functions, colour, categoryName) {
    const functionMap = {};
    const dropdownOptions = [];

    functions.forEach(func => {
        functionMap[func.name] = {
            params: func.params,
            inline: func.inline ?? true
        };
        dropdownOptions.push([func.name, func.name]);
    });

    Blockly.Blocks[blockType] = {
        init: function() {
            this.firstInput_ = this.appendDummyInput('MAIN')
                .appendField('λ')
                .appendField(new Blockly.FieldDropdown(
                    dropdownOptions,
                    this.onFunctionChange.bind(this)
                ), 'NAME');

            this.functionMap_ = functionMap;
            this.currentFunction_ = functions.length > 0 ? functions[0].name : null;
            this.argumentCount_ = 0;

            if (this.currentFunction_) {
                const funcConfig = this.functionMap_[this.currentFunction_];
                const paramNames = funcConfig?.params || [];
                this.updateArguments_(paramNames.length, paramNames);
            }

            this.setOutput(true, null);
            this.setColour(colour);
            this.setTooltip('Call ' + categoryName + ' function');
            this.setHelpUrl('');
        },

        onFunctionChange: function(newValue) {
            this.currentFunction_ = newValue;
            const funcConfig = this.functionMap_[newValue];
            const paramNames = funcConfig?.params || [];
            this.updateArguments_(paramNames.length, paramNames);
            return newValue;
        },

        updateArguments_: function(newCount, paramNames) {
            paramNames = paramNames || [];
            const funcConfig = this.functionMap_[this.currentFunction_];
            const connections = [];

            for (let i = 0; i < this.argumentCount_; i++) {
                const input = this.getInput('ARG' + i);
                if (input && input.connection) {
                    connections[i] = input.connection.targetConnection;
                }
            }

            for (let i = 0; i < this.argumentCount_; i++) {
                if (this.getInput('ARG' + i)) {
                    this.removeInput('ARG' + i);
                }
            }

            if (newCount === 1) {
                const label = paramNames[0] || 'arg';
                this.appendValueInput('ARG0')
                    .appendField(label);
                this.setInputsInline(true);
            } else if (newCount > 1) {
                for (let i = 0; i < newCount; i++) {
                    const label = paramNames[i] || ('arg' + (i + 1));
                    const input = this.appendValueInput('ARG' + i);
                    if (!funcConfig.inline) {
                        input.setAlign(Blockly.ALIGN_RIGHT);
                    }
                    input.appendField(label);
                }
                this.setInputsInline(funcConfig.inline);
            } else {
                this.setInputsInline(false);
            }

            for (let i = 0; i < Math.min(newCount, connections.length); i++) {
                if (connections[i]) {
                    const input = this.getInput('ARG' + i);
                    if (input && input.connection && !input.connection.isConnected()) {
                        input.connection.connect(connections[i]);
                    }
                }
            }

            this.argumentCount_ = newCount;
        },

        mutationToDom: function() {
            const container = Blockly.utils.xml.createElement('mutation');
            container.setAttribute('arguments', this.argumentCount_);
            if (this.currentFunction_) {
                container.setAttribute('function', this.currentFunction_);
            }
            return container;
        },

        domToMutation: function(xmlElement) {
            const argCount = parseInt(xmlElement.getAttribute('arguments'), 10) || 0;
            const functionName = xmlElement.getAttribute('function');

            if (functionName) {
                this.currentFunction_ = functionName;
                this.setFieldValue(functionName, 'NAME');
                const funcConfig = this.functionMap_[functionName];
                const paramNames = funcConfig?.params || [];
                this.updateArguments_(argCount, paramNames);
            } else {
                this.updateArguments_(argCount, []);
            }
        }
    };
}

// register all custom tatu blocks
function registerBlocks() {
    // math functions
    createFunctionCallBlock(
        'math_call_tatu',
        [
            { name: 'math:pi', params: [] },
            { name: 'math:e', params: [] },
            { name: 'math:abs', params: ['x'] },
            { name: 'math:sqrt', params: ['x'] },
            { name: 'math:pow', params: ['x', 'y'] },
            { name: 'math:sin', params: ['x'] },
            { name: 'math:cos', params: ['x'] },
            { name: 'math:tan', params: ['x'] },
            { name: 'math:floor', params: ['x'] },
            { name: 'math:ceil', params: ['x'] },
            { name: 'math:round', params: ['x'] },
            { name: 'math:log', params: ['x'] },
            { name: 'math:exp', params: ['x'] },
            { name: 'math:min', params: ['a', 'b'] },
            { name: 'math:max', params: ['a', 'b'] },
            { name: 'math:rand', params: ['min', 'max'] },
            { name: 'math:between', params: ['x', 'min', 'max'] }
        ],
        230,
        'Math'
    );

    // string functions
    createFunctionCallBlock(
        'string_call_tatu',
        [
            { name: 'str:len', params: ['str'] },
            { name: 'str:upper', params: ['str'] },
            { name: 'str:lower', params: ['str'] },
            { name: 'str:trim', params: ['str'] },
            { name: 'str:reverse', params: ['str'] },
            { name: 'str:contains', params: ['str', 'sub'], inline: false },
            { name: 'str:index', params: ['str', 'sub'], inline: false },
            { name: 'str:starts', params: ['str', 'prefix'], inline: false },
            { name: 'str:ends', params: ['str', 'suffix'], inline: false },
            { name: 'str:repeat', params: ['str', 'num'], inline: false },
            { name: 'str:concat', params: ['str1', 'str2'], inline: false },
            { name: 'str:split', params: ['str', 'delim'], inline: false },
            { name: 'str:join', params: ['vec', 'delim'], inline: false },
            { name: 'str:slice', params: ['str', 'start', 'end'], inline: false },
            { name: 'str:replace', params: ['str', 'old', 'new'], inline: false }
        ],
        160,
        'String'
    );

    // map functions
    createFunctionCallBlock(
        'map_call_tatu',
        [
            { name: 'map:len', params: ['map'] },
            { name: 'map:keys', params: ['map'] },
            { name: 'map:values', params: ['map'] },
            { name: 'map:get', params: ['map', 'key'] },
            { name: 'map:delete', params: ['map', 'key'] },
            { name: 'map:has', params: ['map', 'key'] },
            { name: 'map:merge', params: ['map1', 'map2'] },
            { name: 'map:set', params: ['map', 'key', 'val'] }
        ],
        290,
        'Map'
    );

    // JSON functions
    createFunctionCallBlock(
        'json_call_tatu',
        [
            { name: 'json:encode', params: ['val'] },
            { name: 'json:decode', params: ['str'] }
        ],
        210,
        'JSON'
    );

    // type checking functions
    createFunctionCallBlock(
        'type_check_call_tatu',
        [
            { name: 'is-bool', params: ['val'] },
            { name: 'is-number', params: ['val'] },
            { name: 'is-int', params: ['val'] },
            { name: 'is-string', params: ['val'] },
            { name: 'is-vector', params: ['val'] },
            { name: 'is-map', params: ['val'] },
            { name: 'is-nil', params: ['val'] },
            { name: 'is-function', params: ['val'] }
        ],
        210,
        'Type Check'
    );

    // vector functions
    createFunctionCallBlock(
        'vector_call_tatu',
        [
            { name: 'vec:len', params: ['vec'] },
            { name: 'vec:reverse', params: ['vec'] },
            { name: 'vec:sort', params: ['vec'] },
            { name: 'vec:pop', params: ['vec'] },
            { name: 'vec:get', params: ['vec', 'idx'] },
            { name: 'vec:delete', params: ['vec', 'idx'] },
            { name: 'vec:push', params: ['vec', 'val'] },
            { name: 'vec:contains', params: ['vec', 'val'] },
            { name: 'vec:find', params: ['vec', 'val'] },
            { name: 'vec:concat', params: ['vec1', 'vec2'] },
            { name: 'vec:set', params: ['vec', 'idx', 'val'], inline: false },
            { name: 'vec:slice', params: ['vec', 'start', 'end'], inline: false }
        ],
        260,
        'Vector'
    );

    // time functions
    createFunctionCallBlock(
        'time_call_tatu',
        [
            { name: 'time:now', params: [] },
            { name: 'time:unix', params: ['time'] },
            { name: 'time:year', params: ['time'] },
            { name: 'time:month', params: ['time'] },
            { name: 'time:day', params: ['time'] },
            { name: 'time:hour', params: ['time'] },
            { name: 'time:minute', params: ['time'] },
            { name: 'time:second', params: ['time'] },
            { name: 'time:format', params: ['time', 'layout'], inline: false },
            { name: 'time:parse', params: ['layout', 'str'], inline: false },
            { name: 'time:add', params: ['time', 'duration'], inline: false },
            { name: 'time:sub', params: ['time', 'duration'], inline: false },
            { name: 'time:diff', params: ['time1', 'time2'], inline: false }
        ],
        180,
        'Time'
    );

    // file system functions
    createFunctionCallBlock(
        'filesystem_call_tatu',
        [
            { name: 'fs:read', params: ['path'] },
            { name: 'fs:exists', params: ['path'] },
            { name: 'fs:list', params: ['path'] },
            { name: 'fs:mkdir', params: ['path'] },
            { name: 'fs:delete', params: ['path'] },
            { name: 'fs:size', params: ['path'] },
            { name: 'fs:basename', params: ['path'] },
            { name: 'fs:write', params: ['path', 'content'], inline: false },
            { name: 'fs:append', params: ['path', 'content'], inline: false },
            { name: 'fs:move', params: ['src', 'dest'], inline: false }
        ],
        270,
        'File System'
    );

    // regex functions
    createFunctionCallBlock(
        'regex_call_tatu',
        [
            { name: 'regex:matches', params: ['pattern', 'str'], inline: false },
            { name: 'regex:find', params: ['pattern', 'str'], inline: false },
            { name: 'regex:replace', params: ['pattern', 'str', 'repl'], inline: false }
        ],
        220,
        'Regex'
    );

    // casting functions
    createFunctionCallBlock(
        'casting_call_tatu',
        [
            { name: 'to-string', params: ['val'] },
            { name: 'to-number', params: ['val'] },
            { name: 'to-bool', params: ['val'] }
        ],
        280,
        'Casting'
    );

    // custom function_call_container_tatu block
    Blockly.Blocks['function_call_container_tatu'] = {
        init: function() {
            this.setColour(290);
            this.appendDummyInput()
                .appendField('arguments');
            this.appendStatementInput('STACK');
            this.setTooltip('');
            this.contextMenu = false;
        }
    };

    // custom function_call_item_tatu block
    Blockly.Blocks['function_call_item_tatu'] = {
        init: function() {
            this.setColour(290);
            this.appendDummyInput()
                .appendField('arg');
            this.setPreviousStatement(true);
            this.setNextStatement(true);
            this.setTooltip('');
            this.contextMenu = false;
        }
    };

    // custom function_call_tatu block
    Blockly.Blocks['function_call_tatu'] = {
        init: function() {
            this.setHelpUrl('');
            this.setColour(290);
            this.itemCount_ = 0;
            this.appendDummyInput('FUNCTION_NAME')
                .appendField('λ')
                .appendField(new Blockly.FieldTextInput('name'), 'NAME');
            this.updateShape_();
            this.setOutput(true);
            this.setMutator(new Blockly.icons.MutatorIcon(['function_call_item_tatu'], this));
            this.setTooltip('Call a custom function with arguments.');
        },

        mutationToDom: function() {
            const container = Blockly.utils.xml.createElement('mutation');
            container.setAttribute('items', this.itemCount_);
            return container;
        },

        domToMutation: function(xmlElement) {
            this.itemCount_ = parseInt(xmlElement.getAttribute('items'), 10) || 0;
            this.updateShape_();
        },

        decompose: function(workspace) {
            const containerBlock = workspace.newBlock('function_call_container_tatu');
            containerBlock.initSvg();
            let connection = containerBlock.getInput('STACK').connection;

            for (let i = 0; i < this.itemCount_; i++) {
                const itemBlock = workspace.newBlock('function_call_item_tatu');
                itemBlock.initSvg();
                connection.connect(itemBlock.previousConnection);
                connection = itemBlock.nextConnection;
            }

            return containerBlock;
        },

        compose: function(containerBlock) {
            let itemBlock = containerBlock.getInputTargetBlock('STACK');
            const connections = [];

            while (itemBlock && !itemBlock.isInsertionMarker()) {
                connections.push(itemBlock.valueConnection_);
                itemBlock = itemBlock.nextConnection && itemBlock.nextConnection.targetBlock();
            }

            for (let i = 0; i < this.itemCount_; i++) {
                const connection = this.getInput('ARG' + i);
                if (connection) {
                    const targetConnection = connection.connection.targetConnection;
                    if (targetConnection && connections.indexOf(targetConnection) === -1) {
                        targetConnection.disconnect();
                    }
                }
            }

            this.itemCount_ = connections.length;
            this.updateShape_();

            for (let i = 0; i < this.itemCount_; i++) {
                if (connections[i]) {
                    if (!this.getInput('ARG' + i).connection.targetConnection) {
                        this.getInput('ARG' + i).connection.connect(connections[i]);
                    }
                }
            }
        },

        saveConnections: function(containerBlock) {
            let itemBlock = containerBlock.getInputTargetBlock('STACK');
            let i = 0;

            while (itemBlock) {
                if (itemBlock.isInsertionMarker()) {
                    itemBlock = itemBlock.nextConnection && itemBlock.nextConnection.targetBlock();
                    continue;
                }

                const input = this.getInput('ARG' + i);
                itemBlock.valueConnection_ = input && input.connection.targetConnection;
                i++;
                itemBlock = itemBlock.nextConnection && itemBlock.nextConnection.targetBlock();
            }
        },

        updateShape_: function() {
            for (let i = 0; this.getInput('ARG' + i); i++) {
                this.removeInput('ARG' + i);
            }

            for (let i = 0; i < this.itemCount_; i++) {
                this.appendValueInput('ARG' + i)
                    .setAlign(Blockly.inputs.Align.RIGHT)
                    .appendField('arg' + i);
            }
        }
    };

    // custom comment_tatu block
    Blockly.Blocks['comment_tatu'] = {
        init: function() {
            this.appendDummyInput()
                .appendField('comment')
                .appendField(new Blockly.FieldTextInput('your comment ...'), 'TEXT');
            this.setPreviousStatement(true, null);
            this.setNextStatement(true, null);
            this.setColour(60);
            this.setTooltip('Add a comment to the generated Tatu code');
            this.setHelpUrl('');
        }
    };

    // custom variables_create_tatu block
    Blockly.Blocks['variables_create_tatu'] = {
        init: function() {
            this.appendValueInput('VALUE')
                .appendField('var')
                .appendField(new Blockly.FieldVariable('item'), 'VAR');
            this.setPreviousStatement(true, null);
            this.setNextStatement(true, null);
            this.setColour(330);
            this.setTooltip('Create a new Tatu variable with initial value (generates var declaration)');
            this.setHelpUrl('');
        }
    };

    // custom expression_tatu block
    Blockly.Blocks['expression_tatu'] = {
        init: function() {
            this.appendValueInput('VALUE');
            this.setPreviousStatement(true, null);
            this.setNextStatement(true, null);
            this.setColour(30);
            this.setTooltip('Evaluate an expression (typically used as last statement to return a value)');
            this.setHelpUrl('');
        }
    };

    // custom map_pair_tatu block
    Blockly.Blocks['map_pair_tatu'] = {
        init: function() {
            this.appendValueInput('KEY')
                .appendField('key');
            this.appendValueInput('VALUE')
                .appendField('val');
            this.setInputsInline(true);
            this.setOutput(true, 'MapPair');
            this.setColour(290);
            this.setTooltip('A key-value pair for a map');
        }
    };

    // custom map_create_with_container_tatu block
    Blockly.Blocks['map_create_with_container_tatu'] = {
        init: function() {
            this.setColour(290);
            this.appendDummyInput()
                .appendField('map');
            this.appendStatementInput('STACK');
            this.setTooltip('');
            this.contextMenu = false;
        }
    };

    // custom map_create_with_item_tatu block
    Blockly.Blocks['map_create_with_item_tatu'] = {
        init: function() {
            this.setColour(290);
            this.appendDummyInput()
                .appendField('pair');
            this.setPreviousStatement(true);
            this.setNextStatement(true);
            this.setTooltip('');
            this.contextMenu = false;
        }
    };

    // custom map_create_with_tatu block
    Blockly.Blocks['map_create_with_tatu'] = {
        init: function() {
            this.setHelpUrl('');
            this.setColour(290);
            this.itemCount_ = 2;
            this.appendDummyInput('LABEL')
                .appendField('create map with');
            this.updateShape_();
            this.setOutput(true);
            this.setMutator(new Blockly.icons.MutatorIcon(['map_create_with_item_tatu'], this));
            this.setTooltip('Create a map with key-value pairs.');
        },

        mutationToDom: function() {
            const container = Blockly.utils.xml.createElement('mutation');
            container.setAttribute('items', this.itemCount_);
            return container;
        },

        domToMutation: function(xmlElement) {
            this.itemCount_ = parseInt(xmlElement.getAttribute('items'), 10);
            this.updateShape_();
        },

        saveExtraState: function() {
            return {
                'itemCount': this.itemCount_,
            };
        },

        loadExtraState: function(state) {
            this.itemCount_ = state['itemCount'];
            this.updateShape_();
        },

        decompose: function(workspace) {
            const containerBlock = workspace.newBlock('map_create_with_container_tatu');
            containerBlock.initSvg();
            let connection = containerBlock.getInput('STACK').connection;
            for (let i = 0; i < this.itemCount_; i++) {
                const itemBlock = workspace.newBlock('map_create_with_item_tatu');
                itemBlock.initSvg();
                connection.connect(itemBlock.previousConnection);
                connection = itemBlock.nextConnection;
            }
            return containerBlock;
        },

        compose: function(containerBlock) {
            let itemBlock = containerBlock.getInputTargetBlock('STACK');
            const pairConnections = [];
            while (itemBlock) {
                if (itemBlock.isInsertionMarker()) {
                    itemBlock = itemBlock.getNextBlock();
                    continue;
                }
                pairConnections.push(itemBlock.pairConnection_);
                itemBlock = itemBlock.getNextBlock();
            }

            for (let i = 0; i < this.itemCount_; i++) {
                const input = this.getInput('PAIR' + i);
                if (input) {
                    const pairConnection = input.connection.targetConnection;
                    if (pairConnection && pairConnections.indexOf(pairConnection) === -1) {
                        pairConnection.disconnect();
                    }
                }
            }

            this.itemCount_ = pairConnections.length;
            this.updateShape_();

            for (let i = 0; i < this.itemCount_; i++) {
                if (pairConnections[i]) {
                    const input = this.getInput('PAIR' + i);
                    if (input) {
                        input.connection.connect(pairConnections[i]);
                    }
                }
            }
        },

        saveConnections: function(containerBlock) {
            let itemBlock = containerBlock.getInputTargetBlock('STACK');
            let i = 0;
            while (itemBlock) {
                if (itemBlock.isInsertionMarker()) {
                    itemBlock = itemBlock.getNextBlock();
                    continue;
                }
                const pairInput = this.getInput('PAIR' + i);
                itemBlock.pairConnection_ = pairInput && pairInput.connection.targetConnection;
                i++;
                itemBlock = itemBlock.getNextBlock();
            }
        },

        updateShape_: function() {
            if (this.getInput('EMPTY')) {
                this.removeInput('EMPTY');
            }

            let i;
            for (i = 0; i < this.itemCount_; i++) {
                if (!this.getInput('PAIR' + i)) {
                    this.appendValueInput('PAIR' + i)
                        .setCheck('MapPair');
                }
            }

            for (; this.getInput('PAIR' + i); i++) {
                this.removeInput('PAIR' + i);
            }
        }
    };

    // custom include_file_tatu block
    Blockly.Blocks['include_file_tatu'] = {
        init: function() {
            this.appendDummyInput()
                .appendField('include')
                .appendField(new Blockly.FieldTextInput('lib.tatu'), 'PATH');
            this.setPreviousStatement(true, null);
            this.setNextStatement(true, null);
            this.setColour(310);
            this.setTooltip('Include an external Tatu file');
            this.setHelpUrl('');
        }
    };
}
