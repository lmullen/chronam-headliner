# Download the model
model = Llama-3.2-3B-Instruct.Q6_K.llamafile
download-model : $(model)
$(model) :
	wget https://huggingface.co/Mozilla/Llama-3.2-3B-Instruct-llamafile/resolve/main/$(model)
	chmod +x $(model)

# Run the model as a OpenAI compatible server
.PHONY : run-llm
run-llm : $(model)
# Setting context size to 0 allows maximum context allowed by the model itself
	./$(model) --server --v2 --ctx-size 0

# Run the app
.PHONY : run-app
run-app : 
	go run cmd/main.go

# Send the app a test URL
.PHONY : test-url
test-url :
	curl -X POST http://localhost:8050/chronamurl -H "Content-Type: application/json" \
	-d '{"url": "https://chroniclingamerica.loc.gov/lccn/sn83045462/1925-01-25/ed-1/seq-2/"}'
