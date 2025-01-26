# Download the model
model = Llama-3.2-3B-Instruct.Q6_K.llamafile
download-model : $(model)
$(model) :
	wget https://huggingface.co/Mozilla/Llama-3.2-3B-Instruct-llamafile/resolve/main/$(model)
	chmod +x $(model)

# Run the model as a OpenAI compatible server
.PHONY : run-llm
run-llm : $(model)
	./$(model) --server --v2
