import uvicorn
import io, os, torch
from PIL import Image
from fastapi import FastAPI, UploadFile, File
from fastapi.responses import JSONResponse
from torchvision import transforms
from torchvision.models import resnet18          # ←  use resnet
# ────────────────────────────────────────────────────────
MODEL_PATH  = "best_gpu_tweak.pth"               # your .pth
CLASS_NAMES = ["sad", "neutral", "happy"]
DEVICE      = torch.device("cuda" if torch.cuda.is_available() else "cpu")

preprocess = transforms.Compose([
    transforms.Resize((224, 224)),
    transforms.ToTensor(),
    transforms.Normalize([0.485, 0.456, 0.406],
                         [0.229, 0.224, 0.225])
])

def load_model(ckpt: str):
    model = resnet18(weights=None)               # no download
    model.fc = torch.nn.Sequential(
        torch.nn.Dropout(p=0.2, inplace=True),
        torch.nn.Linear(model.fc.in_features, len(CLASS_NAMES))
    )
    state = torch.load(ckpt, map_location="cpu")
    model.load_state_dict(state, strict=True)    # now keys match
    model.to(DEVICE).eval()
    print(f"✅ Loaded ResNet‑18 from {ckpt} on {DEVICE}")
    return model

model = load_model(MODEL_PATH)

app = FastAPI(title="Facial Expression Classifier")

@app.post("/predict")
async def predict(file: UploadFile = File(...)):
    if not file.content_type.startswith("image"):
        return JSONResponse(status_code=400,
                            content={"error": "File provided is not an image"})

    img_bytes = await file.read()
    img = Image.open(io.BytesIO(img_bytes)).convert("RGB")
    tensor = preprocess(img).unsqueeze(0).to(DEVICE)

    with torch.no_grad():
        logits = model(tensor)
        probs  = torch.softmax(logits, dim=1)[0]
        idx    = int(torch.argmax(probs))
        conf   = float(probs[idx].item())

    return {"class": CLASS_NAMES[idx], "accuracy": conf}

if __name__ == "__main__":
    uvicorn.run("FEC:app", host="0.0.0.0", port=8000)

