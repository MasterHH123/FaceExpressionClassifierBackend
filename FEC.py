import torch
import torch.nn as nn
import torchvision.transforms as transforms
import matplotlib.pyplot as plt
from tqdm import tqdm
from torch.utils.data import DataLoader, random_split
from torchvision.datasets import MNIST
import os
from fastapi import FastAPI, File, UploadFile
from fastapi.responses import JSONResponse
import uvicorn
from PIL import Image
import numpy as np
import io

class EnhancedNN(nn.Module):
    def __init__(self):
        super(EnhancedNN, self).__init__()
        self.fc1 = nn.Linear(28 * 28, 256)  # Changed to 28x28 to match MNIST dimensions
        self.bn1 = nn.BatchNorm1d(256)
        self.dropout1 = nn.Dropout(0.3)
        self.fc2 = nn.Linear(256, 128)
        self.bn2 = nn.BatchNorm1d(128)
        self.dropout2 = nn.Dropout(0.3)
        self.fc3 = nn.Linear(128, 64)
        self.bn3 = nn.BatchNorm1d(64)
        self.dropout3 = nn.Dropout(0.3)
        self.fc4 = nn.Linear(64, 10)

    def forward(self, x):
        x = x.view(-1, 28 * 28)  # Flatten MNIST images to 28*28
        x = torch.relu(self.bn1(self.fc1(x)))
        x = self.dropout1(x)
        x = torch.relu(self.bn2(self.fc2(x)))
        x = self.dropout2(x)
        x = torch.relu(self.bn3(self.fc3(x)))
        x = self.dropout3(x)
        x = self.fc4(x)
        return x


class FECPredictor:
    def __init__(self, model_path="best_gpu_tweak.pth"):
        self.device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
        self.model = EnhancedNN().to(self.device)
        
        # Load the saved model if it exists
        if os.path.exists(model_path):
            self.model.load_state_dict(torch.load(model_path, map_location=self.device))
            self.model.eval()
            print(f"Model loaded from {model_path}")
        else:
            print(f"No model found at {model_path}. Please train the model first.")
            
    def predict(self, image):
        """
        Predict the facial expression from an uploaded image.
        Input image should be in the format expected by the model.
        """
        # Convert to grayscale if needed
        if image.mode != 'L':
            image = image.convert('L')
            
        # Resize to 28x28 if needed
        if image.size != (28, 28):
            image = image.resize((28, 28))
            
        # Convert to tensor and normalize
        transform = transforms.Compose([
            transforms.ToTensor(),
            transforms.Normalize((0.1307,), (0.3081,))
        ])
        
        input_tensor = transform(image).unsqueeze(0).to(self.device)
        
        # Predict
        with torch.no_grad():
            output = self.model(input_tensor)
            _, predicted = torch.max(output, 1)
            probabilities = torch.nn.functional.softmax(output, dim=1)[0]
            
        return {
            "class": predicted.item(),
            "accuracy": float(probabilities[predicted].item()),
            #"probabilities": {i: float(prob) for i, prob in enumerate(probabilities.tolist())}
        }

app = FastAPI(title="Facial Expression Predictor")

# Initialize predictor
predictor = None

@app.on_event("startup")
async def startup_event():
    global predictor
    predictor = FECPredictor()

@app.post("/predict")
async def predict(file: UploadFile = File(...)):
    # Check if the file is an image
    if not file.content_type.startswith("image/"):
        return JSONResponse(
            status_code=400,
            content={"message": "File provided is not an image."}
        )
    
    # Read and process the image
    contents = await file.read()
    image = Image.open(io.BytesIO(contents))
    
    # Make prediction
    result = predictor.predict(image)
    
    return result

if __name__ == "__main__":
    uvicorn.run("FEC:app", host="0.0.0.0", port=8000, reload=True)
