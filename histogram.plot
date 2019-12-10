width=1.0
hist(x,width)=width*floor(x/width)
set boxwidth width*0.95
set style fill solid 0.5
set term svg
set yrange [0:*]
set output "histogram.svg"
plot "test.txt" u (hist($1,width)):(1.0) smooth freq w boxes lc rgb"red" notitle
